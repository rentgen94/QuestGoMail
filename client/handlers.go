package client

import (
	"github.com/rentgen94/QuestGoMail/entities"
	"html/template"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles(
	"./templates/header.html",
	"./templates/footer.html",
	"./templates/index.html",
	"./templates/login.html",
	"./templates/register.html",
	"./templates/game.html",
))

var validPath = regexp.MustCompile("(^/(login|register)/([a-zA-Z0-9]+)$)|(^/)")

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, user *entities.Player) {
	err := templates.ExecuteTemplate(w, tmpl, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func indexHandle(w http.ResponseWriter, r *http.Request) {
	var user *entities.Player
	renderTemplate(w, "index", user)
	return
}

func loginHandle(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login", nil)
}

func registerHandle(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "register", nil)
}

func gameHandle(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "game", nil)
}
