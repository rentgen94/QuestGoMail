package main

import (
	"github.com/rentgen94/QuestGoMail/server"

	"log"
	"net/http"
)

func main() {
	router := server.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
