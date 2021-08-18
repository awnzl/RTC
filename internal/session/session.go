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
			// will I lost previous unread messages? when it will be unread?
			// todo find answers
			msgType, msg, err := s.conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			// msgType should be binary type, since we will use serialized chat-message
			// read message
			// find type
			// find room_id
			// use hub to find or create room
			// fill content to the room

			// s.hub.AddMessage()

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
