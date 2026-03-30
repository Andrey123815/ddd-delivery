package events

import (
	"github.com/google/uuid"
)

const OrderCompletedDomainEventName = "order.completed.event"

type OrderCompletedDomainEvent struct {
	id uuid.UUID
	name string
	
	orderId uuid.UUID
}

func NewOrderCompletedDomainEvent(orderId uuid.UUID) *OrderCompletedDomainEvent {
	return &OrderCompletedDomainEvent{
		id: uuid.New(),
		name: OrderCompletedDomainEventName,
		
		orderId: orderId,
	}
}

func (e *OrderCompletedDomainEvent) GetID() uuid.UUID {
	return e.id
}

func (e *OrderCompletedDomainEvent) GetName() string {
	return e.name
}

func (e *OrderCompletedDomainEvent) GetOrderId() uuid.UUID {
	return e.orderId
}
