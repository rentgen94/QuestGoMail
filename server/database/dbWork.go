package database

import (
	_ "github.com/lib/pq"

	"database/sql"
	"fmt"
	"log"
)

const (
	PostgresUser = "postgres"
	Password     = "qwerty456"
	Server       = "localhost:5432"
	DbName       = "testdb"
	Schema       = "QuestGoMail"
)

var db *sql.DB

func init() {
	DB, err := sql.Open("postgres",
		"postgres://"+PostgresUser+":"+Password+"@"+Server+"/"+DbName+"?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db = DB

	setSCH, err := db.Query("SET SCHEMA '" + Schema + "';")
	if err != nil {
		log.Fatal(err)
	}
	defer setSCH.Close()

	for setSCH.Next() {
		fmt.Println(setSCH.Scan())
	}

	ShowAllUser()
}
