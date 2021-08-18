package domain

import "errors"

type Message interface {
	Content() string
	Destination() string
	Populate() error
	Type() int
}

type ChatMessage struct {
	//type [join, send, leave]
	//time
	//destination room
	//content [destination room (for join  or leave), message text]
}

func (m *ChatMessage) Content() string {
	return "not implemented"
}

func (m *ChatMessage) Destination() string {
	return "not implemented"
}

func (m *ChatMessage) Populate() error {
	return errors.New("not implemented")
}

func (m *ChatMessage) Type() int {
	return -1
}
