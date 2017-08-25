package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/rentgen94/QuestGoMail/server"
	"log"
	"net/http"
	"os"
)

type Configuration struct {
	Port       string  `json:"port"`
	ClientPort string  `json:"clientPort"`
	DriverName string  `json:"driverName"`
	UserName   string  `json:"userName"`
	Password   string  `json:"password"`
	Server     string  `json:"server"`
	DbName     string  `json:"dbName"`
}

func main() {
	conf := getConfig("resources/config/config.json")

	var db = getDb(&conf)
	var env = server.NewEnv(db, 1, 10, 10)
	var routes = server.GetRoutes(env)
	router := server.NewRouter(routes)
	http.ListenAndServe(conf.Port, router)
}

func getDb(conf *Configuration) *sql.DB {
	db, err := sql.Open(conf.DriverName,
		conf.DriverName+"://"+conf.UserName+":"+conf.Password+"@"+conf.Server+"/"+conf.DbName+"?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func getConfig(configPath string) Configuration {
	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		panic(err)
	}
	return configuration
}
