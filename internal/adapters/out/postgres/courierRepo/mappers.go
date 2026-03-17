package courierRepo

import (
	"delivery/internal/core/domain/models/courier"
	"delivery/internal/core/domain/models/kernel"
)

func DomainToDTO(aggregate *courier.Courier) CourierDTO {
	var courierDTO CourierDTO

	courierDTO.Id = aggregate.Id()
	courierDTO.Name = aggregate.Name()
	courierDTO.Speed = aggregate.Speed()

	courierDTO.Location = LocationDTO{
		X: int(aggregate.Location().X()),
		Y: int(aggregate.Location().Y()),
	}

	if aggregate.StoragePlaces() != nil {
		courierDTO.StoragePlaces = make([]StoragePlaceDTO, 0, len(aggregate.StoragePlaces()))
		for _, sp := range aggregate.StoragePlaces() {
			courierDTO.StoragePlaces = append(courierDTO.StoragePlaces, StoragePlaceDTO{
				Id:          sp.Id(),
				CourierID:   aggregate.Id(),
				Name:        sp.Name(),
				TotalVolume: sp.TotalVolume(),
				OrderID:     sp.OrderId(),
			})
		}
	}

	return courierDTO
}

func DtoToDomain(dto CourierDTO) *courier.Courier {
	location, _ := kernel.NewLocation(uint8(dto.Location.X), uint8(dto.Location.Y))
	
	var storagePlaces []*courier.StoragePlace
	if dto.StoragePlaces != nil {
		storagePlaces = make([]*courier.StoragePlace, 0, len(dto.StoragePlaces))
		for _, spDTO := range dto.StoragePlaces {
			sp := courier.RestoreStoragePlace(spDTO.Id, spDTO.Name, spDTO.TotalVolume, spDTO.OrderID)
			storagePlaces = append(storagePlaces, sp)
		}
	}
	
	return courier.RestoreCourierWithStoragePlaces(dto.Id, dto.Name, dto.Speed, location, storagePlaces)
}
