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

func GetRoutes(env *Env) Routes {
	return Routes{
		{
			"GameStartPost",
			http.MethodPost,
			"/game/start/{labyrinth_id}",
			env.GameStartPost,
		},
		{
			"GameQuitPost",
			http.MethodPost,
			"/game/quit",
			env.GameQuitPost,
		},
		{
			"GameCommandPost",
			http.MethodPost,
			"/game/command",
			env.GameCommandPost,
		},

		{
			"PlayerLoginPost",
			http.MethodPost,
			"/player/login",
			env.PlayerLoginPost,
		},

		{
			"PlayerRegisterPost",
			http.MethodPost,
			"/player/register",
			env.PlayerRegisterPost,
		},
		{
			"GameListLabyrinthsGet",
			http.MethodGet,
			"/game/labyrinths",
			env.GameListLabyrinthsGet,
		},
		{
			"GameGetSlotFilling",
			http.MethodGet,
			"/game/slots/{slot_id}",
			env.GameGetSlotFilling,
		},
		{
			"GameLookAroundGet",
			http.MethodGet,
			"/game/look_around",
			env.GameLookAroundGet,
		},
		{
			"GameBagGet",
			http.MethodGet,
			"/game/look_around/bag",
			env.GameBagGet,
		},
		{
			"GameItemsGet",
			http.MethodGet,
			"/game/look_around/entities/items",
			env.GameItemsGet,
		},
		{
			"GameSlotsGet",
			http.MethodGet,
			"/game/look_around/entities/slots",
			env.GameSlotsGet,
		},
		{
			"GameInteractivesGet",
			http.MethodGet,
			"/game/look_around/entities/interactives",
			env.GameInteractivesGet,
		},
		{
			"GameDoorsGet",
			http.MethodGet,
			"/game/look_around/entities/doors",
			env.GameDoorsGet,
		},
	}
}
