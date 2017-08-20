package server

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


var routes = Routes{

	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"GameCommandPost",
		"POST",
		"/game/command",
		GameCommandPost,
	},

	Route{
		"PlayerLoginPost",
		"POST",
		"/player/login",
		PlayerLoginPost,
	},

	Route{
		"PlayerRegisterPost",
		"POST",
		"/player/register",
		PlayerRegisterPost,
	},
}
