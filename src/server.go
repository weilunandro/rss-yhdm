package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Println("Missing url argument")
		return
	}

	index := r.URL.Query().Get("index")
	indexInt, err := strconv.Atoi(index)
	if err != nil {
		indexInt = 1
	}

	feeds := parseBangumi(id, indexInt)
	fmt.Println(feeds)
	resData, err := feeds.ToAtom()
	if err != nil {
		fmt.Errorf(err.Error(), "get Error ")
	}

	w.Header().Add("Content-Type", "application/rss+xml")
	w.Write([]byte(resData))
}
