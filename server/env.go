package server

import (
	"github.com/rentgen94/QuestGoMail/server/database"
)

type Env struct {
	PlayerDAO database.PlayerDAO
}

func NewEnv() Env {
	return Env {
		PlayerDAO: database.NewDBPlayerDAO(database.Init()),
	}
}

