package domain

type ChatRoom struct {
}

type Room interface {
	PushMessage(msg string)
}

func (r *ChatRoom) PushMessage(msg string) {

}

func New() Room {
	return &ChatRoom{}
}
