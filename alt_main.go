package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rentgen94/QuestGoMail/database/dao"
)

func main() {
	var db, err = sql.Open(
		"postgres",
		getDBStr("go_user", "go", "quest_db"),
	)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var labDao = dao.NewLabyrinthDAO(db)
	var lab, labErr = labDao.GetById(1)
	fmt.Println(lab, labErr)

	var inter = lab.StartRoom().Interactives()[1]
	var result, interErr = inter.Interact([]string{"alpha", "beta"}, nil)

	fmt.Println(result, interErr)

}

func getDBStr(user string, pass string, name string) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, pass, name)
}
