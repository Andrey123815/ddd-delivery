package courierRepo

import (
	"delivery/internal/core/domain/models/courier"
	"delivery/internal/core/domain/models/kernel"
)

func DomainToDTO(aggregate *courier.Courier) CourierDTO {
	var courierDTO CourierDTO

	courierDTO.ID = aggregate.ID()
	courierDTO.Name = aggregate.Name()
	courierDTO.Speed = aggregate.Speed().Value()

	courierDTO.Location = LocationDTO{
		X: int(aggregate.Location().X()),
		Y: int(aggregate.Location().Y()),
	}

	if aggregate.StoragePlaces() != nil {
		courierDTO.StoragePlaces = make([]StoragePlaceDTO, 0, len(aggregate.StoragePlaces()))
		for _, sp := range aggregate.StoragePlaces() {
			courierDTO.StoragePlaces = append(courierDTO.StoragePlaces, StoragePlaceDTO{
				ID:          sp.Id(),
				CourierID:   aggregate.ID(),
				Name:        sp.Name(),
				TotalVolume: sp.TotalVolume().Value(),
				OrderID:     sp.OrderId(),
			})
		}
	}

	return courierDTO
}

func DtoToDomain(dto CourierDTO) *courier.Courier {
	location, _ := kernel.NewLocation(uint8(dto.Location.X), uint8(dto.Location.Y))
	speed, _ := courier.NewSpeed(dto.Speed)
	
	var storagePlaces []*courier.StoragePlace
	if dto.StoragePlaces != nil {
		storagePlaces = make([]*courier.StoragePlace, 0, len(dto.StoragePlaces))
		for _, spDTO := range dto.StoragePlaces {
			volume, _ := courier.NewVolume(spDTO.TotalVolume)
			sp := courier.RestoreStoragePlace(spDTO.ID, spDTO.Name, volume, spDTO.OrderID)
			storagePlaces = append(storagePlaces, sp)
		}
	}
	
	return courier.RestoreCourierWithStoragePlaces(dto.ID, dto.Name, speed, location, storagePlaces)
}
