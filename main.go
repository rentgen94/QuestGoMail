package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/rentgen94/QuestGoMail/client"
	"github.com/rentgen94/QuestGoMail/server"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

type Configuration struct {
	Port       string `json:"port"`
	ClientPort string `json:"clientPort"`
	DriverName string `json:"driverName"`
	UserName   string `json:"userName"`
	Password   string `json:"password"`
	Server     string `json:"server"`
	DbName     string `json:"dbName"`
}

func main() {
	conf := getConfig("resources/config/config.json")

	var db = getDb(&conf)
	var env = server.NewEnv(db, 1, 10, 10)
	var routes = server.GetRoutes(env)
	router := server.NewRouter(routes)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	var clientRoutes = client.GetRoutes()
	var clientRouter = client.NewRouter(clientRoutes)
	clientRouter.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))


	c := cors.AllowAll()

	handler := c.Handler(clientRouter)
	go http.ListenAndServe(conf.ClientPort, handler)

	c1 := cors.AllowAll()

	handler1 := c1.Handler(router)
	http.ListenAndServe(conf.Port, handler1)
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
