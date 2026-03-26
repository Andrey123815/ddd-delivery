package order

import (
	"delivery/internal/core/domain/kernel"
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
	return nil
}
