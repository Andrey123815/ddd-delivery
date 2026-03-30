package http

import (
	createCourier "delivery/internal/core/application/usecases/commands/create_courier"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) CreateCourier(c echo.Context) error {
	var body struct {
		Name  string `json:"name"`
		Speed int    `json:"speed"`
	}

	if err := c.Bind(&body); err != nil {
		return err
	}

	command, err := createCourier.NewCreateCourierCommand(body.Name, body.Speed)
	if err != nil {
		return err
	}

	if err := s.createCourierHandler.Handle(c.Request().Context(), command); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, nil)
}
