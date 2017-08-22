package database

import "github.com/rentgen94/QuestGoMail/entities"

type PlayerDAO interface {
	CreatePlayer(player *entities.Player) error
	FindPlayer(player *entities.Player) (*entities.Player, error)
	FindPlayerById(id int) (*entities.Player, error)
	SelectAllPlayers() ([]*entities.Player, error)
}
