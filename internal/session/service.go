package session

import (
	"github.com/awnzl/RTC/internal/hub"
	"github.com/gorilla/websocket"
)

type Service struct {
	hub      *hub.Hub
	sessions map[string]Session
}

func NewService(h *hub.Hub) *Service {
	return &Service{
		hub:      h,
		sessions: make(map[string]Session),
	}
}

func (s *Service) CreateSession(conn *websocket.Conn) Session {
	ses := NewSession(conn, s.hub)
	s.sessions[ses.ID()] = ses

	return ses
}
