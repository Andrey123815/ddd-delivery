package events

import (
	"github.com/google/uuid"
)

const OrderAssignedDomainEventName = "order.assigned.event"

type OrderAssignedDomainEvent struct {
	ID      uuid.UUID
	Name    string
	OrderID uuid.UUID
}

func NewOrderAssignedDomainEvent(orderId uuid.UUID) *OrderAssignedDomainEvent {
	return &OrderAssignedDomainEvent{
		ID:      uuid.New(),
		Name:    OrderAssignedDomainEventName,
		OrderID: orderId,
	}
}

func (e *OrderAssignedDomainEvent) GetID() uuid.UUID {
	return e.ID
}

func (e *OrderAssignedDomainEvent) GetName() string {
	if e.Name != "" {
		return e.Name
	}
	return OrderAssignedDomainEventName
}

func (e *OrderAssignedDomainEvent) GetOrderId() uuid.UUID {
	return e.OrderID
}
