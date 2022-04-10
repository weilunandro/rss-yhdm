package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gorilla/feeds"
)

const BangumiPrefix = "https://www.yhdmp.cc/showp/"

type Bangumi struct {
	id         string
	title      string
	link       string
	updateTime time.Time
	episodes   []*Episode
}

type Episode struct {
	index      int
	id         string
	title      string
	url        string
	updateTime time.Time
}

func generateBangumiUrl(id string) string {
	return BangumiPrefix + id + ".html"
}

func parseBangumi(id string, playlistIdx int) *feeds.Feed {
	url := generateBangumiUrl(id)
	bangumi := crawlBangumi(url, playlistIdx, id)
	episodeIds, _ := GetLocalEpisodes(id)

	newEpisodes := []Episode{}
	for _, episode := range bangumi.episodes {
		if updateTime, ok := episodeIds[episode.id]; ok {
			episode.updateTime = updateTime
		} else {
			episode.updateTime = time.Now()
			newEpisodes = append(newEpisodes, *episode)
		}
	}

	var latestEpisode time.Time
	for _, episode := range bangumi.episodes {
		if episode.updateTime.After(latestEpisode) {
			latestEpisode = episode.updateTime
		}
	}
	bangumi.updateTime = latestEpisode

	if len(newEpisodes) > 0 {
		SaveEpisodes(id, newEpisodes)
	}
	feed := generateFeeds(bangumi)
	return feed
}

func generateFeeds(bangumi *Bangumi) *feeds.Feed {
	result := &feeds.Feed{Title: bangumi.title, Link: &feeds.Link{Href: bangumi.link}}
	var items []*feeds.Item
	for _, episode := range bangumi.episodes {
		item := feeds.Item{Title: episode.title, Link: &feeds.Link{Href: episode.url}, Created: time.Now(), Id: episode.id, Updated: episode.updateTime}
		items = append(items, &item)
	}
	result.Items = items
	result.Id = bangumi.id
	result.Updated = bangumi.updateTime
	return result
}

func crawlBangumi(url string, playlistIdx int, id string) *Bangumi {
	c := colly.NewCollector()
	bangumi := &Bangumi{link: url, id: id}
	bangumi.episodes = []*Episode{}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("requesting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		//fmt.Println("receive response", r)
	})

	// On every a element which has href attribute call callback
	c.OnHTML("div#main0 > div.movurl:nth-child("+strconv.Itoa(playlistIdx)+") > ul", func(e *colly.HTMLElement) {
		e.ForEach("li > a", func(i int, element *colly.HTMLElement) {
			href := element.Attr("href")
			title := element.Attr("title")
			bangumi.episodes = append(bangumi.episodes, &Episode{
				index: i,
				id:    id + "_" + strconv.Itoa(i+1),
				title: title,
				url:   e.Request.AbsoluteURL(href),
			})
		})
	})

	c.OnHTML("div.area > div.fire > div.rate > h1", func(element *colly.HTMLElement) {
		bangumi.title = element.Text
	})

	c.Visit(url)
	return bangumi
}
