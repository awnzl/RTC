package hub

import (
	"github.com/awnzl/RTC/internal/domain"
)

type Hub struct {
	rooms map[string]*domain.Room
}

func New() *Hub {
	return &Hub{
		rooms: make(map[string]*domain.Room),
	}
}

func (h *Hub) PushMessage(msg domain.Message) {

}
