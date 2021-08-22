package domain

import (
	"log"
	"time"
)

type Room interface {
	Subscribe(string, chan<- Message)
	Unsubscribe(string)
	Publish(string, time.Time)
}

type ChatRoom struct {
	roomID      string
	subscribers map[string]chan<- Message
}

func NewRoom(name string) Room {
	return &ChatRoom{
		roomID:      name,
		subscribers: make(map[string]chan<- Message),
	}
}

func (r *ChatRoom) Subscribe(sessionID string, c chan<- Message) {
	// todo: add inform message about joining the room
	r.subscribers[sessionID] = c
}

func (r *ChatRoom) Unsubscribe(sessionID string) {
	// todo: add inform message about leaving the room
	if _, ok := r.subscribers[sessionID]; !ok {
		log.Println("incorrect session id:", sessionID)
		return
	}

	delete(r.subscribers, sessionID)
}

func (r *ChatRoom) Publish(msg string, t time.Time) {
	for _, c := range r.subscribers {
		go func(c chan<- Message) {
			c <- Message{
				Content:  msg,
				Room: r.roomID,
				Time:     t,
				Type:     Send,
			}
		}(c)
	}
}
