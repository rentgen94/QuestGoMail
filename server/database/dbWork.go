package database

import (
	_ "github.com/lib/pq"

	"database/sql"
	"log"
)

const (
	PostgresUser = "postgres"
	Password     = "qwerty456"
	Server       = "localhost:5432"
	DbName       = "testdb"
	Schema       = "QuestGoMail"
)

func Init() *sql.DB {
	db, err := sql.Open("postgres",
		"postgres://"+PostgresUser+":"+Password+"@"+Server+"/"+DbName+"?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	setSCH, err := db.Query("SET SCHEMA '" + Schema + "';")
	if err != nil {
		log.Fatal(err)
	}
	defer setSCH.Close()

	for setSCH.Next() {
	}

	return db
}
