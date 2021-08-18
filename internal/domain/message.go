package domain

import (
	"time"
)

type MessageType int

const (
	join MessageType = iota
	leave
	send
)

type Message struct {
	Content     string
	Destination string
	Time        time.Time
	Type        MessageType
}
