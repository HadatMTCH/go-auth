package template

import (
	"base-api/infra/context/service"
	"github.com/labstack/echo/v4"
)

type Template interface {
	RegistrationUser(c echo.Context) error
	Login(c echo.Context) error
	Profile(c echo.Context) error
}

func New(serviceCtx *service.ServiceContext) Template {
	return &template{
		ServiceContext: serviceCtx,
	}
}
