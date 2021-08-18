package session

import (
	"log"

	"github.com/gorilla/websocket"

	"github.com/awnzl/RTC/internal/hub"
)

type Session struct {
	conn *websocket.Conn
	hub  *hub.Hub
	stop chan struct{}
}

func New(conn *websocket.Conn, h *hub.Hub) *Session {
	s := &Session{
		conn: conn,
		hub:  h,
		stop: make(chan struct{}),
	}

	return s
}

func (s *Session) Run() {
	go s.reader()
}

func (s *Session) Stop() {
	close(s.stop)
	s.conn.Close()
}

func (s *Session) reader() {
	for {
		select {
		case <-s.stop:
			log.Println("Session is stopped")
			return
		default:
			msgType, msg, err := s.conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			// todo: if got "close message", use s.Stop()

			err = s.conn.WriteMessage(msgType, msg)
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}

func (s *Session) writer() {
	log.Println("writer not implemented yet")
}
