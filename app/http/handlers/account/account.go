package account

import (
	"base-api/constants"
	"base-api/data/models"
	"base-api/infra/context/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type accountHandler struct {
	repo *repository.RepositoryContext
}

func (h *accountHandler) Register(c echo.Context) error {
	var req models.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Remark: "invalid request"})
	}

	exists, err := h.repo.AccountRepository.CheckExistingCustomer(req.Nik, req.NoHp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Remark: "internal error"})
	}
	if exists {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Remark: "NIK or phone already registered"})
	}

	noRekening := uuid.New().String()
	account := &models.Account{
		Nama: req.Nama,
		Nik:  req.Nik,
		NoHp: req.NoHp,
	}

	if err := h.repo.AccountRepository.CreateAccount(noRekening, account); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Remark: "internal error"})
	}

	return c.JSON(http.StatusOK, models.RegisterResponse{NoRekening: noRekening})
}

func (h *accountHandler) Tabung(c echo.Context) error {
	var req models.TabungTarikRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Remark: "invalid request"})
	}

	saldo, err := h.repo.AccountRepository.UpdateSaldo(req.NoRekening, req.Nominal, true)
	if err != nil {
		if err == constants.ErrAccountNotFound {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Remark: "account not found"})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Remark: "internal error"})
	}

	return c.JSON(http.StatusOK, models.SaldoResponse{Saldo: saldo})
}

func (h *accountHandler) Tarik(c echo.Context) error {
	var req models.TabungTarikRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Remark: "invalid request"})
	}

	saldo, err := h.repo.AccountRepository.UpdateSaldo(req.NoRekening, req.Nominal, false)
	if err != nil {
		if err == constants.ErrAccountNotFound {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Remark: "account not found"})
		} else if err == constants.ErrInsufficientBalance {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Remark: "insufficient balance"})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Remark: "internal error"})
	}

	return c.JSON(http.StatusOK, models.SaldoResponse{Saldo: saldo})
}

func (h *accountHandler) GetSaldo(c echo.Context) error {
	noRekening := c.Param("no_rekening")
	saldo, err := h.repo.AccountRepository.GetSaldo(noRekening)
	if err != nil {
		if err == constants.ErrAccountNotFound {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Remark: "account not found"})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Remark: "internal error"})
	}
	return c.JSON(http.StatusOK, models.SaldoResponse{Saldo: saldo})
}
