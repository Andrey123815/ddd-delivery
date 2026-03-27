package http

import (
	getNoCompletedOrders "delivery/internal/core/application/usecases/queries/get_no_completed_orders"
	"net/http"

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

	return c.JSON(http.StatusOK, result)
}
