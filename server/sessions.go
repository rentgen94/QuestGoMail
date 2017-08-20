package server

import (
	"github.com/gorilla/sessions"
)

const (
	sessionName = "Session-for-new-user"
	authToken   = "auth_token"
)

var store = sessions.NewCookieStore([]byte("server-cookie-store"))
