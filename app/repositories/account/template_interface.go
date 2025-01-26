package account

import (
	"base-api/data/models"
	"base-api/infra/db"
)

type Account interface {
	CheckExistingCustomer(nik, noHp string) (bool, error)
	CreateAccount(noRekening string, account *models.Account) error
	UpdateSaldo(noRekening string, nominal int64, isTabung bool) (int64, error)
	GetSaldo(noRekening string) (int64, error)
}

func New(db *db.DB) Account {
	return &AccountRepository{
		db.Master(),
	}
}
