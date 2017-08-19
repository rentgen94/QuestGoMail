package main

import (
	_ "github.com/lib/pq"
	"html/template"
	"net/http"
	"regexp"
	"fmt"
	"strconv"

	"github.com/rentgen94/QuestGoMail/services"
)

var users map[int]*services.User

var templates = template.Must(template.ParseFiles(
	"templates/header.html",
	"templates/footer.html",
	"templates/index.html",
	"templates/login.html",
	"templates/register.html",
))

var validPath = regexp.MustCompile("(^/(login|register)/([a-zA-Z0-9]+)$)|(^/)")

func renderTemplate(w http.ResponseWriter, tmpl string, user *services.User) {
	err := templates.ExecuteTemplate(w, tmpl, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func makeHandler(fn func (http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

func indexHandel(w http.ResponseWriter, r *http.Request) {
	var user *services.User
	for _, value := range users {
		if value != nil {
			user = value
			break
		}
	}
	renderTemplate(w, "index", user)
	return
}

func loginHandel(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("uname")
	password := r.FormValue("psw")

	if password != "" && name != "" {
		founded := services.FindUser(w, r, name, password)
		if founded != nil {
			users[founded.Id] = &services.User{founded.Id,founded.Name, founded.Password}
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	renderTemplate(w, "login", nil)
}

func logoutHandle(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id != "" {
		toDel, err := strconv.Atoi(id);
		if err != nil {
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}
		delete(users, toDel)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	fmt.Println(users)
}

func registerHandel(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("uname")
	password := r.FormValue("psw")

	if password != "" && name != "" {
		registered := services.FindUser(w, r, name, password)
		if registered != nil {
			renderTemplate(w, "register", registered)
			return
		}
		services.CreateUser(w, r, name, password)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	renderTemplate(w, "register", nil)
}

func main() {
	users = make(map[int]*services.User, 0)

	http.HandleFunc("/", makeHandler(indexHandel))
	http.HandleFunc("/login", makeHandler(loginHandel))
	http.HandleFunc("/logout", makeHandler(logoutHandle))
	http.HandleFunc("/register", makeHandler(registerHandel))

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	http.ListenAndServe(":8080", nil)
}