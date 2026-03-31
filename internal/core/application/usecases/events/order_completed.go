package events

import (
	"github.com/google/uuid"
)

const OrderCompletedDomainEventName = "order.completed.event"

type OrderCompletedDomainEvent struct {
	ID      uuid.UUID
	Name    string
	OrderID uuid.UUID
}

func NewOrderCompletedDomainEvent(orderId uuid.UUID) *OrderCompletedDomainEvent {
	return &OrderCompletedDomainEvent{
		ID:      uuid.New(),
		Name:    OrderCompletedDomainEventName,
		OrderID: orderId,
	}
}

func (e *OrderCompletedDomainEvent) GetID() uuid.UUID {
	return e.ID
}

func (e *OrderCompletedDomainEvent) GetName() string {
	if e.Name != "" {
		return e.Name
	}
	return OrderCompletedDomainEventName
}

func (e *OrderCompletedDomainEvent) GetOrderId() uuid.UUID {
	return e.OrderID
}
