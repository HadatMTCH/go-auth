package health

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type health struct{}

func (h *health) Check(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
