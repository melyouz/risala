package internal

type Message struct {
	QueueName string `json:"queueName" validate:"required"`
	Payload   string `json:"payload" validate:"required"`
}
