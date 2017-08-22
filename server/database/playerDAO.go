package database

import (
	"database/sql"
	"log"

	"github.com/rentgen94/QuestGoMail/entities"
)

const (
	RegisterOk        = "\"Player login successful.\""
	AlreadyRegistered = "\"Player already registered.\""
	RegisterError     = "\"Player register error.\""
	PlayerNotFound    = "\"Player not found.\""
	PlayerFoundError  = "\"Player found error.\""
	PlayerFoundOk     = "\"Player found successful.\""

	createPlayer   = "INSERT INTO users(name, password) VALUES ($1, $2);"
	findPlayer     = "SELECT * FROM users WHERE name=$1 AND password=$2"
	findById       = "SELECT * FROM users WHERE id=$1"
	findAllPlayers = "SELECT * FROM users"
)

type dbPlayerDAO struct {
	db *sql.DB
}

func NewDBPlayerDAO(db *sql.DB) PlayerDAO {
	var result = new(dbPlayerDAO)
	result.db = db
	return result
}

func (dao *dbPlayerDAO) CreatePlayer(player *entities.Player) string {
	founded, _ := dao.FindPlayer(player)
	if founded == nil {
		_, err := dao.db.Exec(createPlayer, player.Login, player.Password)
		if err != nil {
			return RegisterError
		}
	} else {
		return AlreadyRegistered
	}

	return RegisterOk
}

func (dao *dbPlayerDAO) FindPlayer(player *entities.Player) (*entities.Player, string) {
	row := dao.db.QueryRow(findPlayer, player.Login, player.Password)

	founded := new(entities.Player)
	err := row.Scan(&founded.Id, &founded.Login, &founded.Password)
	if err == sql.ErrNoRows {
		return nil, PlayerNotFound
	} else if err != nil {
		return nil, PlayerFoundError
	}

	return founded, PlayerFoundOk
}

func (dao *dbPlayerDAO) FindPlayerById(id int) (*entities.Player, string) {
	row := dao.db.QueryRow(findById, id)

	founded := new(entities.Player)
	err := row.Scan(&founded.Id, &founded.Login, &founded.Password)
	if err == sql.ErrNoRows {
		return nil, PlayerNotFound
	} else if err != nil {
		return nil, PlayerFoundError
	}

	return founded, PlayerFoundOk
}

func (dao *dbPlayerDAO) SelectAllPlayers() []*entities.Player {
	rows, err := dao.db.Query(findAllPlayers)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	usrs := make([]*entities.Player, 0)
	for rows.Next() {
		us := new(entities.Player)
		err := rows.Scan(&us.Id, &us.Login, &us.Password)
		if err != nil {
			log.Fatal(err)
		}
		usrs = append(usrs, us)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return usrs
}
