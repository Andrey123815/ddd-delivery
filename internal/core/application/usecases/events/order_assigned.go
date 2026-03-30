package events

import (
	"github.com/google/uuid"
)

const OrderAssignedDomainEventName = "order.assigned.event"

type OrderAssignedDomainEvent struct {
	id uuid.UUID
	name string

	orderId uuid.UUID
}

func NewOrderAssignedDomainEvent(orderId uuid.UUID) *OrderAssignedDomainEvent {
	return &OrderAssignedDomainEvent{
		id: uuid.New(),
		name: OrderAssignedDomainEventName,
		
		orderId: orderId,
	}
}

func (e *OrderAssignedDomainEvent) GetID() uuid.UUID {
	return e.id
}

func (e *OrderAssignedDomainEvent) GetName() string {
	return e.name
}

func (e *OrderAssignedDomainEvent) GetOrderId() uuid.UUID {
	return e.orderId
}
