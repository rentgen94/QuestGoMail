package database

import "github.com/rentgen94/QuestGoMail/entities"

type PlayerDAO interface {
	CreatePlayer(player *entities.Player) string
	FindPlayer(player *entities.Player) (*entities.Player, string)
	FindPlayerById(id int) (*entities.Player, string)
	SelectAllPlayers() []*entities.Player
}
