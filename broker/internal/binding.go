package internal

import "github.com/google/uuid"

type Binding struct {
	Id         uuid.UUID `json:"id"`
	Queue      string    `json:"queue" validate:"required"`
	RoutingKey string    `json:"routingKey"`
}
