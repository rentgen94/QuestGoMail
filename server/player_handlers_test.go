package server

import (
	"github.com/gorilla/sessions"
	"github.com/rentgen94/QuestGoMail/management"
	"github.com/rentgen94/QuestGoMail/server/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEnc_PlayerRegisterPost_Successful(t *testing.T) {
	var env = &Env{PlayerDAO: new(mocks.ExistPlayerDAOMock)}
	var inputMsg = "{\"login\":\"qqq\", \"password\":\"111\"}"

	rr := getRecorder(t, "POST", "/player/register",
		strings.NewReader(inputMsg), env.PlayerRegisterPost)

	checkStatus(t, http.StatusCreated, rr)

	checkBody(t, RegisterOk, rr)
}

func TestEnv_PlayerRegisterPost_Fail(t *testing.T) {
	var env = &Env{PlayerDAO: new(mocks.NotExistPlayerDAOMock)}
	var inputMsg = "{\"login\":\"qqq\", \"password\":\"111\"}"

	rr := getRecorder(t, "POST", "/player/register",
		strings.NewReader(inputMsg), env.PlayerRegisterPost)

	checkStatus(t, http.StatusConflict, rr)

	checkBody(t, RegisterError, rr)
}

func TestEnv_PlayerLoginPost_Successful(t *testing.T) {
	var env = &Env{PlayerDAO: new(mocks.ExistPlayerDAOMock),
		Store:      sessions.NewCookieStore([]byte("server-cookie-store")),
		playerId:   "player_id",
		cookieName: "quest_go_mail",
		gameId:     "game_id",
		currGameId:    1,
		Pool:       management.NewManagerPool(1, 10, 10),
	}
	var inputMsg = "{\"login\":\"qqq\", \"password\":\"111\"}"

	rr := getRecorder(t, "POST", "/player/login",
		strings.NewReader(inputMsg), env.PlayerLoginPost)

	checkStatus(t, http.StatusOK, rr)

	checkBody(t, PlayerFoundOk, rr)
}

func TestEnv_PlayerLoginPost_Fail(t *testing.T) {
	var env = &Env{PlayerDAO: new(mocks.NotExistPlayerDAOMock),
		Store:      sessions.NewCookieStore([]byte("server-cookie-store")),
		playerId:   "player_id",
		cookieName: "quest_go_mail",
	}
	var inputMsg = "{\"login\":\"qqq\", \"password\":\"1112132131\"}"

	rr := getRecorder(t, "POST", "/player/login",
		strings.NewReader(inputMsg), env.PlayerLoginPost)

	checkStatus(t, http.StatusNotFound, rr)

	checkBody(t, PlayerNotFound, rr)
}

func getRecorder(t *testing.T, method string, url string, body io.Reader,
	handlerFunc http.HandlerFunc) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(rr, req)

	return rr
}

func checkStatus(t *testing.T, expStatus int, rr *httptest.ResponseRecorder) {
	if status := rr.Code; status != expStatus {
		t.Errorf("handler returned wrong status code: got %v <-| want |-> %v",
			status, expStatus)
	}
}

func checkBody(t *testing.T, expMsg string, rr *httptest.ResponseRecorder) {
	if rr.Body.String() != expMsg {
		t.Errorf("handler returned unexpected body: got %v <-| want |-> %v",
			rr.Body.String(), expMsg)
	}
}
