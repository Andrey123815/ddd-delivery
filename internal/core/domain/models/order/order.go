package order

import (
	"delivery/internal/core/application/usecases/events"
	"delivery/internal/core/domain/models/kernel"
	"delivery/internal/pkg/ddd"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Order struct {
	id uuid.UUID
	courierId uuid.UUID
	location kernel.Location
	volume int
	status OrderStatus

	domainEvents []ddd.DomainEvent
}

func NewOrder(location kernel.Location, volume int) (*Order, error) {
	if location.IsEmpty() {
		return nil, errors.New("Попытка создать order с пустым location")
	}
	if volume <= 0 {
		return nil, fmt.Errorf("Попытка создать товар с некорректным volume: %d", volume)
	}

	return &Order{
		id:        uuid.New(),
		courierId: uuid.Nil,
		location:  location,
		volume:    volume,
		status:    OrderStatusCreated,
	}, nil
}

func RestoreOrder(id uuid.UUID, courierId uuid.UUID, location kernel.Location, volume int, status OrderStatus) *Order {
	return &Order{
		id:        id,
		courierId: courierId,
		location:  location,
		volume:    volume,
		status:    status,
	}
}

func (o *Order)Equals(other Order) bool {
	return o.id == other.id;
}

func (o *Order)Id() uuid.UUID {
	return o.id;
}

func (o *Order)CourierId() uuid.UUID {
	return o.courierId;
}

func (o *Order)Location() kernel.Location {
	return o.location;
}

func (o *Order)Volume() int {
	return o.volume;
}

func (o *Order)Status() OrderStatus {
	return o.status;
}

func (o *Order)Assign(courierId uuid.UUID) error {
	if o.status == OrderStatusAssigned {
		return errors.New("Попытка назначить уже назначенный заказ")
	}
	if o.status == OrderStatusCompleted {
		return errors.New("Попытка назначить уже завершенный заказ")
	}
	if courierId == uuid.Nil {
		return errors.New("Попытка назначить заказ курьеру с пустым идентификатором")
	}

	o.courierId = courierId
	o.status = OrderStatusAssigned

	o.RaiseDomainEvent(events.NewOrderAssignedDomainEvent(o.id))

	return nil
}

func (o *Order)Complete(courierId uuid.UUID) error {
	if o.status == OrderStatusCreated {
		return errors.New("Попытка завершить неназначенный заказ")
	}
	if o.status == OrderStatusCompleted {
		return errors.New("Попытка завершить уже завершенный заказ")
	}
	if o.courierId != courierId {
		return errors.New("Попытка завершить заказ другим курьером")
	}

	o.status = OrderStatusCompleted

	o.RaiseDomainEvent(events.NewOrderCompletedDomainEvent(o.id))

	return nil
}

func (o *Order)GetDomainEvents() []ddd.DomainEvent {
	return o.domainEvents
}

func (o *Order)ClearDomainEvents() {
	o.domainEvents = []ddd.DomainEvent{}
}

func (o *Order)RaiseDomainEvent(event ddd.DomainEvent) {
	o.domainEvents = append(o.domainEvents, event)
}
