package courierRepo

import (
	"delivery/internal/core/domain/models/order"

	"github.com/google/uuid"
)

type OrderDTO struct {
	ID        uuid.UUID   `gorm:"type:uuid;primaryKey"`
	Volume    int
	Status    order.OrderStatus `gorm:"type:varchar(20)"`
	Location  LocationDTO `gorm:"embedded;embeddedPrefix:location_"`
	CourierID uuid.UUID   `gorm:"type:uuid;index"`
}

type LocationDTO struct {
	X int
	Y int
}

func (OrderDTO) TableName() string {
	return "orders"
}
