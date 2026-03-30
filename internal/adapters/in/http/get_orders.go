package http

import (
	getNoCompletedOrders "delivery/internal/core/application/usecases/queries/get_no_completed_orders"
	"delivery/internal/generated/servers"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/labstack/echo/v4"
)

func (s *Server) GetOrders(c echo.Context) error {
	query, err := getNoCompletedOrders.NewGetNoCompletedOrdersQuery()
	if err != nil {
		return err
	}

	result, err := s.getNoCompletedOrdersHandler.Handle(c.Request().Context(), query)
	if err != nil {
		return err
	}

	out := make([]servers.Order, 0, len(result))
	for _, row := range result {
		out = append(out, toGetOrdersHTTPResponse(row))
	}
	return c.JSON(http.StatusOK, out)
}

func toGetOrdersHTTPResponse(row getNoCompletedOrders.GetNoCompletedOrdersResponse) servers.Order {
	return servers.Order{
		Id:       openapi_types.UUID(row.Id),
		Location: servers.Location{X: row.Location.X, Y: row.Location.Y},
	}
}
