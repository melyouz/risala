package internal

type Queue struct {
	Name       string         `json:"name" validate:"required"`
	Durability DurabilityType `json:"durability" validate:"required,oneof=durable transient"`
	Messages   []Message      `json:"-" validate:"dive"`
}

func (q *Queue) PublishMessage(message Message) {
	q.Messages = append(q.Messages, message)
}
