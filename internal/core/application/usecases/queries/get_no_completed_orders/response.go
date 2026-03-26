package getNoCompletedOrders

import (
	"delivery/internal/core/domain/models/order"

	"github.com/google/uuid"
)

type GetNoCompletedOrdersResponse struct {
	Id        uuid.UUID `json:"id"`
	CourierId uuid.UUID `json:"courierId,omitempty"`
	Location  LocationDTO `json:"location"`
	Volume    int `json:"volume"`
	Status    string `json:"status"`
}

type LocationDTO struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func ToResponse(o *order.Order) GetNoCompletedOrdersResponse {
	response := GetNoCompletedOrdersResponse{
		Id:     o.Id(),
		Volume: o.Volume(),
		Status: string(o.Status()),
		Location: LocationDTO{
			X: int(o.Location().X()),
			Y: int(o.Location().Y()),
		},
	}

	if o.CourierId() != uuid.Nil {
		response.CourierId = o.CourierId()
	}

	return response
}
