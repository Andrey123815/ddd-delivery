package courierRepo

import (
	"github.com/google/uuid"
)

type CourierDTO struct {
	Id        uuid.UUID   `gorm:"type:uuid;primaryKey"`
	Name      string      `gorm:"type:varchar(255)"`
	Speed     int         `gorm:"type:int"`
	Location  LocationDTO `gorm:"embedded;embeddedPrefix:location_"`
	StoragePlaces []StoragePlaceDTO `gorm:"foreignKey:CourierID"`
}

type StoragePlaceDTO struct {
	Id        uuid.UUID   `gorm:"type:uuid;primaryKey"`
	CourierID uuid.UUID   `gorm:"type:uuid;index"`
	Name      string      `gorm:"type:varchar(255)"`
	TotalVolume int       `gorm:"type:int"`
	OrderID uuid.UUID   `gorm:"type:uuid;index"`
}

type LocationDTO struct {
	X int
	Y int
}

func (CourierDTO) TableName() string {
	return "couriers"
}

func (StoragePlaceDTO) TableName() string {
	return "storage_places"
}
