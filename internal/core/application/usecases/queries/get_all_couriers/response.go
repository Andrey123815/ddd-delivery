package getAllCouriers

import (
	"delivery/internal/core/domain/models/courier"

	"github.com/google/uuid"
)

type GetAllCouriersResponse struct {
	Id            uuid.UUID          `json:"id"`
	Name          string             `json:"name"`
	Speed         int                `json:"speed"`
	Location      LocationDTO        `json:"location"`
	StoragePlaces []StoragePlaceDTO  `json:"storagePlaces"`
}

type LocationDTO struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type StoragePlaceDTO struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	TotalVolume int       `json:"totalVolume"`
	IsOccupied  bool      `json:"isOccupied"`
	OrderId     uuid.UUID `json:"orderId,omitempty"`
}

func ToResponse(c *courier.Courier) GetAllCouriersResponse {
	response := GetAllCouriersResponse{
		Id:    c.Id(),
		Name:  c.Name(),
		Speed: c.Speed().Value(),
		Location: LocationDTO{
			X: int(c.Location().X()),
			Y: int(c.Location().Y()),
		},
		StoragePlaces: make([]StoragePlaceDTO, 0),
	}

	if c.StoragePlaces() != nil {
		for _, sp := range c.StoragePlaces() {
			spDTO := StoragePlaceDTO{
				Id:          sp.Id(),
				Name:        sp.Name(),
				TotalVolume: sp.TotalVolume().Value(),
				IsOccupied:  sp.IsOccupied(),
			}
			if sp.IsOccupied() {
				spDTO.OrderId = sp.OrderId()
			}
			response.StoragePlaces = append(response.StoragePlaces, spDTO)
		}
	}

	return response
}
