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
		GameLookAroundGet,
	},
	Route{
		"GameBagGet",
		"GET",
		"/game/look_around/bag",
		GameBagGet,
	},
	Route{
		"GameItemsGet",
		"GET",
		"/game/look_around/entities/items",
		GameItemsGet,
	},
	Route{
		"GameSlotsGet",
		"GET",
		"/game/look_around/entities/slots",
		GameSlotsGet,
	},
	Route{
		"GameInteractivesGet",
		"GET",
		"/game/look_around/entities/interactives",
		GameIteractivesGet,
	},
	Route{
		"GameDoorsGet",
		"GET",
		"/game/look_around/entities/doors",
		GameDoorsGet,
	},

}
