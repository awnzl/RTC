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
	id   string
	conn *websocket.Conn
	hub  *hub.Hub
	stop chan struct{}
	send chan domain.Message // register this as subscriber's incoming channel
}

func NewSession(conn *websocket.Conn, h *hub.Hub) Session {
	s := &ChatSession{
		id:   uuid.New().String(),
		conn: conn,
		hub:  h,
		stop: make(chan struct{}), //todo: were I use this?
		send: make(chan domain.Message),
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
	// close(s.send)
	s.conn.Close()
}

// todo: remove struct
type msgtmp struct {
	Type    int    `json:"type"`
	Time    string `json:"time"`
	Room    string `json:"room"`
	Content string `json:"content"`
}

func (s *ChatSession) reader() {
	for {
		select {
		case <-s.stop:
			log.Println("Session is stopped")
			return
		default:
			// will I lost previous unread messages? when it will be unread?
			// todo find answers
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

			switch msg.Type {
			case domain.Join:
				log.Println("[session] reader: join", msg.Room)
				s.hub.JoinRoom(s.id, msg.Room, s.send)
				// find or create room, join it, send message "Client Name joined"
			case domain.Leave:
				log.Println("[session] reader: leave")
				s.hub.LeaveRoom(s.id, msg.Room)
				close(s.send)
				// find room, leave it, send message "Client Name left room"
				// destroy link between client and room
			case domain.Send:
				log.Println("[session] reader: send")
				s.hub.PushMessage(msg.Room, msg.Content, msg.Time)
				// find room, send message
			default:
				log.Panic("incorrect domain message type", wsMsgType)
			}
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
