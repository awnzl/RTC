package hub

import (
	"errors"
	"time"

	"github.com/awnzl/RTC/internal/domain"
)

var errIncorrectRoomID = errors.New("incorrect room identifier")

type Hub struct {
	rooms map[string]domain.Room
}

func New() *Hub {
	return &Hub{
		rooms: make(map[string]domain.Room),
	}
}

func (h *Hub) JoinRoom(sessionID, name string, send chan<- domain.Message) {
	r, ok := h.rooms[name]
	if !ok {
		r = domain.NewRoom(name)
		h.rooms[name] = r
	}

	r.Subscribe(sessionID, send)
}

func (h *Hub) LeaveRoom(sessionID, name string) error {
	r, ok := h.rooms[name]
	if !ok {
		return errIncorrectRoomID
	}

	r.Unsubscribe(sessionID)

	return nil
}

func (h *Hub) PushMessage(name, msg string, t time.Time) error {
	r, ok := h.rooms[name]
	if !ok {
		return errIncorrectRoomID
	}

	r.Publish(msg, t)

	return nil
}
