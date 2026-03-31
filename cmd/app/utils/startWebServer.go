package utils

import (
	"delivery/cmd"
	"delivery/internal/adapters/in/http/problems"
	"delivery/internal/generated/servers"
	"fmt"
	"log"
	"net/http"

	httpAdapter "delivery/internal/adapters/in/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func StartWebServer(cr *cmd.CompositionRoot, port string) {
	httpServer, err := httpAdapter.NewServer(
		cr.NewCreateCourierHandler(),
		cr.NewCreateOrderHandler(),
		cr.NewGetAllCouriersHandler(),
		cr.NewGetNoCompletedOrdersHandler(),
	)
	if err != nil {
		log.Fatalf("cannot create http server: %v", err)
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))
	e.HTTPErrorHandler = problems.EchoErrorHandler

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Healthy")
	})

	servers.RegisterHandlers(e, httpServer)

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%s", port)))
}
