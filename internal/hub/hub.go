package hub

import (
	"log"
	"time"

	"github.com/awnzl/RTC/internal/domain"
)

type Hub struct {
	rooms map[string]domain.Room
}

func New() *Hub {
	return &Hub{
		rooms: make(map[string]domain.Room),
	}
}

func (h *Hub) JoinRoom(sessionID, name string, send chan<- domain.Message) {
	log.Println("hub join room:", name)

	r, ok := h.rooms[name]
	if !ok {
		r = domain.NewRoom(name)
		h.rooms[name] = r
	}

	r.Subscribe(sessionID, send)
}

func (h *Hub) LeaveRoom(sessionID, name string) {
	r, ok := h.rooms[name]
	if !ok {
		log.Println("incorrect room identifier:", name)
		return
	}

	r.Unsubscribe(sessionID)
}

func (h *Hub) PushMessage(name, msg string, t time.Time) {
	log.Println("[hub] PushMessage:", name, msg)

	r, ok := h.rooms[name]
	if !ok {
		log.Println("incorrect room identifier:", name)
		return
	}

	r.Publish(msg, t)
}
