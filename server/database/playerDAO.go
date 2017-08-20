package database

import (
	"log"
	"fmt"
	"database/sql"

	"github.com/rentgen94/QuestGoMail/entities"
)

const (
	RegisterOk        = "Пользователь успешно зарегистрирован."
	AlreadyRegistered = "Пользователь уже присутствует в базе данных."
	RegisterError     = "Не удалось зарегестрировать пользователя."
	PlayerNotFound    = "Не удалось найти игрока."
	PlayerFoundError  = "Ошибка поиска игрока при обращении к базе данных."
	PlayerFoundOk     = "Игрок найден."
)

func CreateUser(player *entities.Player) string {
	founded, _ := FindPlayer(player)
	if founded == nil {
		_, err := db.Exec("INSERT INTO users(name, password) VALUES ($1, $2);", player.Login, player.Password)
		if err != nil {
			return RegisterError
		}
	} else {
		return AlreadyRegistered
	}

	return RegisterOk
}

func FindPlayer(player *entities.Player) (*entities.Player, string) {
	row := db.QueryRow("SELECT * FROM users WHERE name=$1 AND password=$2", player.Login, player.Password)

	founded := new(entities.Player)
	err := row.Scan(&founded.Id, &founded.Login, &founded.Password)
	if err == sql.ErrNoRows {
		return nil, PlayerNotFound
	} else if err != nil {
		return nil, PlayerFoundError
	}

	return founded, PlayerFoundOk
}

func FindPlayerById(id int) (*entities.Player, string) {
	row := db.QueryRow("SELECT * FROM users WHERE id=$1", id)

	founded := new(entities.Player)
	err := row.Scan(&founded.Id, &founded.Login, &founded.Password)
	if err == sql.ErrNoRows {
		return nil, PlayerNotFound
	} else if err != nil {
		return nil, PlayerFoundError
	}

	return founded, PlayerFoundOk
}

func ShowAllUser() []*entities.Player {
	rows, err := db.Query("SELECT * FROM users")
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

	for _, us := range usrs {
		fmt.Printf("%d, %s, %s\n", us.Id, us.Login, us.Password)
	}

	return usrs
}