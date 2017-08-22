package server

import (
	"net/http"
)

var env = NewEnv()

type Route struct {
	Name               string
	Method             string
	Pattern            string
	HandlerFunc        http.HandlerFunc
}

type Routes []Route

var routes = Routes{

	Route{
		"Index",
		"GET",
		"/",
		env.Index,
	},

	Route{
		"GameCommandPost",
		"POST",
		"/game/command",
		env.GameCommandPost,
	},

	Route{
		"PlayerLoginPost",
		"POST",
		"/player/login",
		env.PlayerLoginPost,
	},

	Route{
		"PlayerRegisterPost",
		"POST",
		"/player/register",
		env.PlayerRegisterPost,
	},
}