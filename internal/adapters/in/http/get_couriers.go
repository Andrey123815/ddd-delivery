package http

import (
	getAllCouriers "delivery/internal/core/application/usecases/queries/get_all_couriers"
	"delivery/internal/generated/servers"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/labstack/echo/v4"
)

func (s *Server) GetCouriers(c echo.Context) error {
	query, err := getAllCouriers.NewGetAllCouriersQuery()
	if err != nil {
		return err
	}

	result, err := s.getAllCouriersHandler.Handle(c.Request().Context(), query)
	if err != nil {
		return err
	}

	out := make([]servers.Courier, 0, len(result))
	for _, row := range result {
		out = append(out, toGetCouriersHTTPResponse(row))
	}
	return c.JSON(http.StatusOK, out)
}

func toGetCouriersHTTPResponse(row getAllCouriers.GetAllCouriersResponse) servers.Courier {
	return servers.Courier{
		Id:       openapi_types.UUID(row.Id),
		Name:     row.Name,
		Location: servers.Location{X: row.Location.X, Y: row.Location.Y},
	}
}
