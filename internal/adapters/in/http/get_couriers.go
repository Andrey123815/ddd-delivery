package http

import (
	getAllCouriers "delivery/internal/core/application/usecases/queries/get_all_couriers"
	"net/http"

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

	return c.JSON(http.StatusOK, result)
}
