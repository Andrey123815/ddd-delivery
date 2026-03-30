package http

import (
	createOrder "delivery/internal/core/application/usecases/commands/create_order"
	"delivery/internal/generated/queues/basketeventspb"
	"delivery/internal/generated/servers"
	"net/http"

	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (s *Server) CreateOrder(c echo.Context) error {
	command, err := createOrder.NewCreateOrderCommand(
		"http-test-basket",
		&basketeventspb.Address{Street: "Тестировочная"},
		[]*basketeventspb.Item{{Title: "test", Quantity: 1}},
		&basketeventspb.DeliveryPeriod{From: 9, To: 18},
		1,
	)
	if err != nil {
		return err
	}

	orderID, err := s.createOrderHandler.Handle(c.Request().Context(), command)
	if err != nil {
		return err
	}

	oid := openapi_types.UUID(orderID)
	return c.JSON(http.StatusCreated, servers.CreateOrderResponse{OrderId: &oid})
}
