package account

import (
	"base-api/infra/context/repository"
	"github.com/labstack/echo/v4"
)

type Account interface {
	Register(c echo.Context) error
	Tabung(c echo.Context) error
	Tarik(c echo.Context) error
	GetSaldo(c echo.Context) error
}

func New(repo *repository.RepositoryContext) Account {
	return &accountHandler{repo: repo}
}
