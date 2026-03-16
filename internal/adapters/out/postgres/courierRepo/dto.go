package courierRepo

import (
	"delivery/internal/core/domain/models/order"

	"github.com/google/uuid"
)

type CourierDTO struct {
	ID        uuid.UUID   `gorm:"type:uuid;primaryKey"`
	Name      string      `gorm:"type:varchar(255)"`
	Speed     int         `gorm:"type:int"`
	Location  LocationDTO `gorm:"embedded;embeddedPrefix:location_"`
	Volume    int
	Status    order.OrderStatus `gorm:"type:varchar(20)"`
}

type LocationDTO struct {
	X int
	Y int
}

func (CourierDTO) TableName() string {
	return "couriers"
}
