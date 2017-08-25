package dao

import (
	"database/sql"

	"errors"
	"github.com/rentgen94/QuestGoMail/entities"
)

const (
	createPlayer   = "INSERT INTO users(name, password) VALUES ($1, $2);"
	findPlayer     = "SELECT * FROM users WHERE name=$1 AND password=$2"
	playerExists   = "SELECT count(*) FROM users WHERE name=$1"
	findById       = "SELECT * FROM users WHERE id=$1"
	findAllPlayers = "SELECT * FROM users"

	playerAlreadyExists = "Player already exists"
)

type PlayerDAO interface {
	CreatePlayer(player *entities.Player) error
	FindPlayer(player *entities.Player) (*entities.Player, error)
	Exist(player *entities.Player) (bool, error)
	FindPlayerById(id int) (*entities.Player, error)
	SelectAllPlayers() ([]*entities.Player, error)
}

type dbPlayerDAO struct {
	db *sql.DB
}

func NewDBPlayerDAO(db *sql.DB) PlayerDAO {
	var result = new(dbPlayerDAO)
	result.db = db
	return result
}

func (dao *dbPlayerDAO) CreatePlayer(player *entities.Player) error {
	exist, findErr := dao.Exist(player)
	if findErr != nil {
		return findErr
	}

	if exist {
		return errors.New(playerAlreadyExists)
	}

	_, execErr := dao.db.Exec(createPlayer, player.Login, player.Password)
	if execErr != nil {
		return execErr
	}

	return nil
}

func (dao *dbPlayerDAO) FindPlayer(player *entities.Player) (*entities.Player, error) {
	row := dao.db.QueryRow(findPlayer, player.Login, player.Password)

	err := row.Scan(&player.Id, &player.Login, &player.Password)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func (dao *dbPlayerDAO) Exist(player *entities.Player) (bool, error) {
	row := dao.db.QueryRow(playerExists, player.Login)

	var cnt int
	err := row.Scan(&cnt)
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}

func (dao *dbPlayerDAO) FindPlayerById(id int) (*entities.Player, error) {
	row := dao.db.QueryRow(findById, id)

	founded := new(entities.Player)
	err := row.Scan(&founded.Id, &founded.Login, &founded.Password)
	if err == sql.ErrNoRows {
		return nil, err
	}

	return founded, nil
}

func (dao *dbPlayerDAO) SelectAllPlayers() ([]*entities.Player, error) {
	rows, err := dao.db.Query(findAllPlayers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usrs := make([]*entities.Player, 0)
	for rows.Next() {
		us := new(entities.Player)
		err := rows.Scan(&us.Id, &us.Login, &us.Password)
		if err != nil {
			return usrs, err
		}
		usrs = append(usrs, us)
	}

	if err = rows.Err(); err != nil {
		return usrs, err
	}

	return usrs, nil
}
