package account

import (
	"base-api/constants"
	"base-api/data/models"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type AccountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) CheckExistingCustomer(nik, noHp string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM accounts WHERE nik = $1 OR no_hp = $2`
	err := r.db.Get(&count, query, nik, noHp)
	return count > 0, err
}

func (r *AccountRepository) CreateAccount(noRekening string, account *models.Account) error {
	query := `INSERT INTO accounts (no_rekening, nama, nik, no_hp) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, noRekening, account.Nama, account.Nik, account.NoHp)
	return err
}

func (r *AccountRepository) UpdateSaldo(noRekening string, nominal int64, isTabung bool) (int64, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Update saldo
	operation := "+"
	jenis := "tabung"
	if !isTabung {
		var saldo int64
		query := `SELECT saldo FROM accounts WHERE no_rekening = $1`
		err = tx.Get(&saldo, query, noRekening)
		if err != nil {
			if err == sql.ErrNoRows {
				return 0, constants.ErrAccountNotFound
			}
			return 0, err
		}
		if saldo < nominal {
			return 0, constants.ErrInsufficientBalance
		}
		operation = "-"
		jenis = "tarik"
	}

	var saldo int64
	query := `UPDATE accounts SET saldo = saldo ` + operation + ` $1 WHERE no_rekening = $2 RETURNING saldo`
	err = tx.Get(&saldo, query, nominal, noRekening)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, constants.ErrAccountNotFound
		}
		return 0, err
	}

	// Insert transaction
	_, err = tx.Exec(`INSERT INTO transactions (no_rekening, nominal, jenis) VALUES ($1, $2, $3)`, noRekening, nominal, jenis)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return saldo, nil
}

func (r *AccountRepository) GetSaldo(noRekening string) (int64, error) {
	var saldo int64
	query := `SELECT saldo FROM accounts WHERE no_rekening = $1`
	err := r.db.Get(&saldo, query, noRekening)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, constants.ErrAccountNotFound
		}
		return 0, err
	}
	return saldo, nil
}
