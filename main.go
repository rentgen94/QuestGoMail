package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/rentgen94/QuestGoMail/server"
	"log"
	"net/http"
)

const (
	PostgresUser = "postgres"
	Password     = "45274245"
	Server       = "localhost:5432"
	DbName       = "testDB"
	Schema       = "QuestGoMail"
)

func main() {
	var db = getDb()
	var env = server.NewEnv(db, 1, 10, 10)
	var routes = server.GetRoutes(env)
	router := server.NewRouter(routes)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func getDb() *sql.DB {
	db, err := sql.Open("postgres",
		"postgres://"+PostgresUser+":"+Password+"@"+Server+"/"+DbName+"?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
