package courierRepo

import (
	"delivery/internal/core/domain/models/courier"
	"delivery/internal/core/domain/models/kernel"
)

func DomainToDTO(aggregate *courier.Courier) CourierDTO {
	var courierDTO CourierDTO

	courierDTO.ID = aggregate.ID()
	courierDTO.Name = aggregate.Name()
	courierDTO.Speed = aggregate.Speed()

	courierDTO.Location = LocationDTO{
		X: int(aggregate.Location().X()),
		Y: int(aggregate.Location().Y()),
	}

	return courierDTO
}

func DtoToDomain(dto CourierDTO) *courier.Courier {
	var aggregate *courier.Courier

	location, _ := kernel.NewLocation(uint8(dto.Location.X), uint8(dto.Location.Y))
	aggregate = courier.RestoreCourier(dto.ID, dto.Name, dto.Speed, location)
	
	return aggregate
}
