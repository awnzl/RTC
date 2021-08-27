package v1

import (
	"log"
	"net/http"

	"github.com/awnzl/RTC/internal/session"
	"github.com/gorilla/websocket"
)

type API struct {
	sessionManager *session.Service
	upgrader       websocket.Upgrader
}

func New(sm *session.Service) *API {
	api := &API{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		sessionManager: sm,
	}

	api.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	return api
}

func (a *API) RegisterHandlers() {
	http.HandleFunc("/v1/status", a.status)
	http.HandleFunc("/v1/ws", a.serveWS)
}

func (a *API) status(w http.ResponseWriter, r *http.Request) {
	// todo: write response 200, Online
}

func (a *API) serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	err = conn.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}

	wsSession := a.sessionManager.CreateSession(conn)
	wsSession.Run()
}
