package server

import (
	"log"
	"net/http"

	"github.com/awnzl/RTC/internal/hub"
	"github.com/awnzl/RTC/internal/session"
	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader websocket.Upgrader
	server   http.Server
	hub      *hub.Hub
}

func New() *Server {
	s := &Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		server: http.Server{
			Addr:    "localhost:8080",
			Handler: nil,
		},
		hub: hub.New(),
	}

	s.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	http.HandleFunc("/status", s.status)
	http.HandleFunc("/ws", s.serveWS)

	return s
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *Server) status(w http.ResponseWriter, r *http.Request) {
	//write response 200, Online
}

func (s *Server) serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("client connected")

	err = conn.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}

	wsSession := session.New(conn, s.hub) // where to store sessions? server or hub
	wsSession.Run()
	log.Println("after run")
}
