package main

import (
	"log"
	"net/http"
)

func main() {
	InitDb()
	CreateTable()

	addr := ":7171"

	http.HandleFunc("/rss/yhdm", Handler)
	log.Println("listening on ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))

}
