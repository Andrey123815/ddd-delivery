package courier

import (
	"delivery/internal/core/domain/kernel"
	"errors"

	"github.com/google/uuid"

	orderModel "delivery/internal/core/domain/models/order"
)

type Courier struct {
	id uuid.UUID
	name string
	speed int
	location kernel.Location
	storagePlaces []*StoragePlace
}

func NewCourier(name string, speed int, location kernel.Location) (*Courier, error) {
	if name == "" {
		return nil, errors.New("имя курьера не может быть пустым")
	}
	if speed <= 0 {
		return nil, errors.New("скорость курьера должна быть больше нуля")
	}
	if location.IsEmpty() {
		return nil, errors.New("локация курьера задана некорректно")
	}
	return &Courier{
		id:            uuid.New(),
		name:          name,
		speed:         speed,
		location:      location,
		storagePlaces: nil,
	}, nil
}

func (c *Courier)Equals(other Courier) bool {
	return c.id == other.id;
}

func (c *Courier)ID() uuid.UUID {
	return c.id;
}

func (c *Courier)Name() string {
	return c.name;
}

func (c *Courier)Speed() int {
	return c.speed;
}

func (c *Courier)Location() kernel.Location {
	return c.location;
}

func (c *Courier)StoragePlaces() []*StoragePlace  {
	return c.storagePlaces;
}

func (c *Courier)AddStoragePlace(name string, volume int) error {
	storagePlace, err := NewStoragePlace(name, volume)
	if err != nil {
		return err
	}

	if (c.storagePlaces == nil) {
		c.storagePlaces = make([]*StoragePlace, 0)
	}

	c.storagePlaces = append(c.storagePlaces, storagePlace)

	return nil
}

func (c *Courier) CanTakeOrder(order orderModel.Order) (bool, error) {
	for _, storagePlace := range c.storagePlaces {
		canStore, err := storagePlace.CanStore(order.Volume())
		if err != nil {
			return false, err
		}
		if canStore {
			return true, nil
		}
	}
	
	return false, nil
}

func (c *Courier) CompleteOrder(order *orderModel.Order) error {
	if order == nil {
		return errors.New("заказ не может быть nil")
	}
	for _, storagePlace := range c.storagePlaces {
		if storagePlace.OrderId() != order.Id() {
			continue
		}
		if err := storagePlace.Clear(order.Id()); err != nil {
			return err
		}
		if err := order.Complete(c.id); err != nil {
			return err
		}
		return nil
	}
	return errors.New("у курьера нет данного заказа в местах хранения")
} 

func (c *Courier) CalculateTimeToLocation(target kernel.Location) (float64, error) {
	if c.location.IsEmpty() {
		return 0, errors.New("у курьера некорректно заданы текущие координаты")
	}
	if target.IsEmpty() {
		return 0, errors.New("целевая локация задана некорректно")
	}
	if c.speed <= 0 {
		return 0, errors.New("скорость курьера должна быть больше нуля")
	}
	distance, err := c.location.DistanceTo(target)
	if err != nil {
		return 0, err
	}
	return float64(distance) / float64(c.speed), nil
}

func (c *Courier) Move(target kernel.Location) error {
	if target.IsEmpty() {
		return errors.New("попытка переместить курьера по некорректным координатам")
	}
	c.location = target
	return nil
}

func (c *Courier) findStoragePlaceByOrderId(orderId uuid.UUID) (*StoragePlace, error) {
	if c.storagePlaces == nil {
		return nil, errors.New("у данного курьера нет доступных мест для перемещения грузов")
	}
	for _, storagePlace := range c.storagePlaces {
		if !storagePlace.isOccupied() || storagePlace.OrderId() != orderId {
			continue
		}
		return storagePlace, nil
	}
	return nil, errors.New("у данного курьера нет таких заказов в работе")
}
