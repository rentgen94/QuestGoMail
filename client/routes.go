package client

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func GetRoutes() Routes {
	return Routes{
		{
			"Index",
			http.MethodGet,
			"/",
			makeHandler(indexHandle),
		},
		{
			"Login",
			http.MethodGet,
			"/player/login",
			makeHandler(loginHandle),
		},
	}
}
