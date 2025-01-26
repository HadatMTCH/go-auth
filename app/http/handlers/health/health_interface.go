package health

import (
	"github.com/labstack/echo/v4"
)

type Health interface {
	Check(c echo.Context) error
}

func New() Health {
	return &health{}
}
