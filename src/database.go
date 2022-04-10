package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDb() {
	var err error
	config := mysql.Config{
		User:                 "root",
		Passwd:               "Wl@19890919",
		Addr:                 "127.0.0.1:3306",
		Net:                  "tcp",
		DBName:               "rsshub",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	// 准备数据库连接池
	db, err = sql.Open("mysql", config.FormatDSN())
	checkError(err)

	// 设置最大连接数
	db.SetMaxOpenConns(25)
	// 设置最大空闲连接数
	db.SetMaxIdleConns(25)
	// 设置每个链接的过期时间
	db.SetConnMaxLifetime(5 * time.Minute)

	// 尝试连接，失败会报错
	err = db.Ping()
	checkError(err)
}

func CreateTable() {
	createYHDMTale := `create table if not exists yhdm 
	(
    id          int auto_increment primary key ,
    bangumi_id  varchar(32) not null,
    update_time datetime    not null,
    episode_id  varchar(32) not null
	);
	`
	_, err := db.Exec(createYHDMTale)
	checkError(err)
}

func GetLocalEpisodes(id string) (map[string]time.Time, error) {
	query := ("SELECT episode_id, update_time FROM yhdm WHERE bangumi_id = ?")

	rows, err := db.Query(query, id)
	if err != nil {
		log.Fatal(err)
	}

	var result = make(map[string]time.Time)
	for rows.Next() {
		var id string
		var updateTime time.Time
		err = rows.Scan(&id, &updateTime)
		if err != nil {
			log.Fatal(err)
		}
		result[id] = updateTime
	}

	return result, nil
}

func SaveEpisodes(banguiId string, episodeIds []Episode) {
	for _, episode := range episodeIds {
		_, err := saveEpisode(banguiId, episode)
		if err != nil {
			log.Println(err)
		}
	}
}

func saveEpisode(bangumiId string, episode Episode) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO yhdm (bangumi_id, episode_id, update_time) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	dt := episode.updateTime.Format(time.RFC3339)
	rs, err := stmt.Exec(bangumiId, episode.id, dt)
	if err != nil {
		return 0, err
	}

	if id, _ := rs.LastInsertId(); id > 0 {
		return id, nil
	}
	return 0, err
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
