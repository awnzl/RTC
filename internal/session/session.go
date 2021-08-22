package session

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/awnzl/RTC/internal/domain"
	"github.com/awnzl/RTC/internal/hub"
)

type Session interface {
	ID() string
	Run()
	Stop()
}

type ChatSession struct {
	id         string
	activeRoom string
	conn       *websocket.Conn
	hub        *hub.Hub
	stop       chan struct{}
	send       chan domain.Message // register this as subscriber's incoming channel
}

func NewSession(conn *websocket.Conn, h *hub.Hub) Session {
	s := &ChatSession{
		id:   uuid.New().String(),
		conn: conn,
		hub:  h,
		send: make(chan domain.Message),
		stop: make(chan struct{}), //todo: were I use this?
	}

	return s
}

func (s *ChatSession) ID() string {
	return s.id
}

func (s *ChatSession) Run() {
	go s.reader()
	go s.writer()
}

func (s *ChatSession) Stop() {
	close(s.stop)
	s.conn.Close()
}

func (s *ChatSession) reader() {
	for {
		select {
		case <-s.stop:
			log.Println("Session is stopped")
			return

		default:
			wsMsgType, bts, err := s.conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			// todo: do I need this?
			if wsMsgType != websocket.TextMessage {
				continue
			}

			var msg domain.Message
			if json.Unmarshal(bts, &msg) != nil {
				log.Println("[session] reader:", err)
				continue
			}

			s.sendMessage(msg)
		}
	}
}

func (s *ChatSession) writer() {
	for {
		select {
		case <-s.stop:
			log.Println("Session is stopped")
			return

		default:
			bts, err := json.Marshal(<-s.send)
			if err != nil {
				log.Println("session writer:", err)
				continue
			}

			err = s.conn.WriteMessage(1, bts)
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}

func (s *ChatSession) sendMessage(msg domain.Message) {
	switch msg.Type {
	case domain.Join:
		if s.activeRoom == msg.Room {
			return
		}

		log.Println("[session] reader: join", msg.Room)
		s.activeRoom = msg.Room
		s.hub.JoinRoom(s.id, msg.Room, s.send)

	case domain.Leave:
		log.Println("[session] reader: leave")

		err := s.hub.LeaveRoom(s.id, msg.Room)
		if err != nil {
			log.Println(err, msg.Room)
			return
		}

		close(s.send)

	case domain.Send:
		log.Println("[session] reader: send")

		err := s.hub.PushMessage(msg.Room, msg.Content, msg.Time)
		if err != nil {
			log.Println(err, msg.Room)
		}
	}
}
