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
	Route{
		"GameLookAroundGet",
		"GET",
		"/game/look_around",
		env.GameLookAroundGet,
	},
	Route{
		"GameBagGet",
		"GET",
		"/game/look_around/bag",
		env.GameBagGet,
	},
	Route{
		"GameItemsGet",
		"GET",
		"/game/look_around/entities/items",
		env.GameItemsGet,
	},
	Route{
		"GameSlotsGet",
		"GET",
		"/game/look_around/entities/slots",
		env.GameSlotsGet,
	},
	Route{
		"GameInteractivesGet",
		"GET",
		"/game/look_around/entities/interactives",
		env.GameIteractivesGet,
	},
	Route{
		"GameDoorsGet",
		"GET",
		"/game/look_around/entities/doors",
		env.GameDoorsGet,
	},

}
