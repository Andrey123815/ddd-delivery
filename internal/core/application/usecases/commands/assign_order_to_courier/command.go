package assignCourier

import (
	"github.com/google/uuid"
)

type AssignOrderToCourierCommand struct {}


func NewAssignOrderToCourierCommand(orderId uuid.UUID, courierId uuid.UUID) (*AssignOrderToCourierCommand, error) {
	return &AssignOrderToCourierCommand{}, nil
}
