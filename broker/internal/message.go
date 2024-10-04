package internal

type Message struct {
	Payload string `json:"payload" validate:"required"`
}
