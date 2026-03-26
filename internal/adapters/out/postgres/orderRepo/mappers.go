package courierRepo

import (
	"delivery/internal/core/domain/models/kernel"
	"delivery/internal/core/domain/models/order"
)

func DomainToDTO(aggregate *order.Order) OrderDTO {
	var orderDTO OrderDTO

	orderDTO.ID = aggregate.Id()
	orderDTO.CourierID = aggregate.CourierId()
	orderDTO.Volume = aggregate.Volume()
	orderDTO.Status = aggregate.Status()

	orderDTO.Location = LocationDTO{
		X: int(aggregate.Location().X()),
		Y: int(aggregate.Location().Y()),
	}

	return orderDTO
}

func DtoToDomain(dto OrderDTO) *order.Order {
	var aggregate *order.Order

	location, _ := kernel.NewLocation(uint8(dto.Location.X), uint8(dto.Location.Y))
	aggregate = order.RestoreOrder(dto.ID, dto.CourierID, location, dto.Volume, dto.Status)
	
	return aggregate
}
