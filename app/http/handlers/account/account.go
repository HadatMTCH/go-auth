package account

import (
	"base-api/constants"
	"base-api/data/models"
	"base-api/infra/context/repository"
	"base-api/utils"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type accountHandler struct {
	repo *repository.RepositoryContext
}

// Register godoc
// @Summary      Register a new account
// @Description  Register a new customer account with NIK, phone number, and name
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        request body models.RegisterRequest true "Registration details"
// @Success      200  {object}  models.RegisterResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /api/accounts/register [post]
func (h *accountHandler) Register(c echo.Context) error {
	var req models.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Remark: "invalid request"})
	}

	exists, err := h.repo.AccountRepository.CheckExistingCustomer(req.Nik, req.NoHp)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Remark: "internal error"})
	}
	if exists {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Remark: "NIK or phone already registered"})
	}

	randomNumber, _ := utils.GenerateRandomNumber(20)
	noRekening := fmt.Sprint(randomNumber)
	account := &models.Account{
		Nama: req.Nama,
		Nik:  req.Nik,
		NoHp: req.NoHp,
	}
	fmt.Println(noRekening)
	if err := h.repo.AccountRepository.CreateAccount(noRekening, account); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Remark: "internal error"})
	}

	return c.JSON(http.StatusOK, models.RegisterResponse{NoRekening: noRekening})
}

// Tabung godoc
// @Summary      Deposit money
// @Description  Deposit money into an account
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        request body models.TabungTarikRequest true "Deposit details"
// @Success      200  {object}  models.SaldoResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /api/accounts/tabung [post]
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

// Tarik godoc
// @Summary      Withdraw money
// @Description  Withdraw money from an account
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        request body models.TabungTarikRequest true "Withdrawal details"
// @Success      200  {object}  models.SaldoResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /api/accounts/tarik [post]
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

// GetSaldo godoc
// @Summary      Get account balance
// @Description  Retrieve current balance of an account
// @Tags         accounts
// @Produce      json
// @Param        no_rekening path string true "Account number"
// @Success      200  {object}  models.SaldoResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /api/accounts/saldo/{no_rekening} [get]
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
