package problems

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func EchoErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	if he, ok := err.(*echo.HTTPError); ok {
		_ = c.JSON(he.Code, Problem{
			Status:  he.Code,
			Message: http.StatusText(he.Code),
		})
		return
	}

	problem := FromError(err)
	_ = c.JSON(problem.Status, problem)
}
