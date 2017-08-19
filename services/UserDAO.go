package services

import (
	"log"
	"fmt"
	"net/http"
	"database/sql"
)

type User struct {
	Id       int
	Name     string
	Password string
}

func CreateUser(w http.ResponseWriter, r *http.Request, name, pass string) {
	_, err := db.Exec("INSERT INTO users(name, password) VALUES ($1, $2);", name, pass)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	//rowsAffected, err := result.RowsAffected()
	//if err != nil {
	//	http.Error(w, http.StatusText(500), 500)
	//	return
	//}
}

func FindUser(w http.ResponseWriter, r *http.Request, name, pass string) *User {
	row := db.QueryRow("SELECT * FROM users WHERE name=$1 AND password=$2", name, pass)

	founded := new(User)
	err := row.Scan(&founded.Id, &founded.Name, &founded.Password)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return nil
	}

	return founded
}

func ShowAllUser() []*User {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	usrs := make([]*User, 0)
	for rows.Next() {
		us := new(User)
		err := rows.Scan(&us.Id, &us.Name, &us.Password)
		if err != nil {
			log.Fatal(err)
		}
		usrs = append(usrs, us)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	for _, us := range usrs {
		fmt.Printf("%d, %s, %s\n", us.Id, us.Name, us.Password)
	}

	return usrs
}