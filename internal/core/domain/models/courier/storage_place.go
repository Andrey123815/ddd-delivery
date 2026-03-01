package courier

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type StoragePlace struct {
	id uuid.UUID
	name string
	totalVolume int
	orderId uuid.UUID
}

func NewStoragePlace(name string, totalVolume int) (*StoragePlace, error) {
	if name == "" {
		return nil, errors.New("Название места хранения не может быть пустым")
	}

	if totalVolume <= 0 {
		return nil, errors.New("Объем места хранения не может быть меньше или равным нулю")
	}

	return &StoragePlace{
		id:          uuid.New(),
		name:        name,
		totalVolume: totalVolume,
	}, nil
}

func (s *StoragePlace) Equals(other *StoragePlace) bool {
	return s.id == other.id
}

func (s *StoragePlace) Id() uuid.UUID {
	return s.id
}

func (s *StoragePlace) Name() string {
	return s.name
}

func (s *StoragePlace) TotalVolume() int {
	return s.totalVolume
}

func (s *StoragePlace) OrderId() uuid.UUID {
	return s.orderId
}

func (s *StoragePlace) CanStore(volume int) (bool, error) {
	if volume <= 0 {return false, errors.New("Попытка разместить пустой или отрицательный объем") }

	canStore := s.orderId == uuid.Nil && s.totalVolume >= volume

	return canStore, nil
}

func (s *StoragePlace) Store(orderId uuid.UUID, volume int) error {
	canStore, err := s.CanStore(volume)

	if err != nil { return err }

	if canStore == false { return errors.New("Место хранения не может разместить заказ") }

	s.orderId = orderId

	return nil
}

func (s *StoragePlace) Clear(orderId uuid.UUID) error {
	if s.isOccupied() == false { return errors.New("Попытка очистить пустое хранилище") }

	if s.orderId != orderId { return fmt.Errorf("ордера с id %s нет в хранилище", orderId) }

	s.orderId = uuid.Nil

	return nil
}

func (s *StoragePlace) isOccupied() bool {
	return s.orderId != uuid.Nil
}
