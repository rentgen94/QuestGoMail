package server

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/rentgen94/QuestGoMail/server/mocks"
	"strings"
	"github.com/gorilla/sessions"
	"io"
	"github.com/rentgen94/QuestGoMail/management"
)

func TestEnv_Index(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(env.Index)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v <-| want |-> %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `Hello world!`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v <-| want |-> %v",
			rr.Body.String(), expected)
	}
}

func TestEnc_PlayerRegisterPost_Successful(t *testing.T) {
	var env = &Env{PlayerDAO: new(mocks.ExistPlayerDAOMock),}
	var inputMsg = "{\"login\":\"qqq\", \"password\":\"111\"}"

	rr := getRecorder(t, "POST", "/player/register",
		strings.NewReader(inputMsg), env.PlayerRegisterPost)

	checkStatus(t, http.StatusCreated, rr)

	checkBody(t, RegisterOk, rr)
}

func TestEnv_PlayerRegisterPost_Fail(t *testing.T) {
	var env = &Env{PlayerDAO: new(mocks.NotExistPlayerDAOMock),}
	var inputMsg = "{\"login\":\"qqq\", \"password\":\"111\"}"

	rr := getRecorder(t, "POST", "/player/register",
		strings.NewReader(inputMsg), env.PlayerRegisterPost)

	checkStatus(t, http.StatusConflict, rr)

	checkBody(t, RegisterError, rr)
}

func TestEnv_PlayerLoginPost_Successful(t *testing.T) {
	var env = &Env{PlayerDAO: new(mocks.ExistPlayerDAOMock),
		Store: sessions.NewCookieStore([]byte("server-cookie-store")),
		playerId:   "player_id",
		cookieName: "quest_go_mail",
		gameId:     "game_id",
		curGame:    1,
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
		Store: sessions.NewCookieStore([]byte("server-cookie-store")),
		playerId:   "player_id",
		cookieName: "quest_go_mail",
	}
	var inputMsg = "{\"login\":\"qqq\", \"password\":\"1112132131\"}"

	rr := getRecorder(t, "POST", "/player/login",
		strings.NewReader(inputMsg), env.PlayerLoginPost)

	checkStatus(t, http.StatusBadRequest, rr)

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