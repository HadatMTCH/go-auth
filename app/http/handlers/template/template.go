package template

import (
	"base-api/data/models"
	"base-api/infra/context/service"
	"base-api/infra/middleware"
	"base-api/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type template struct {
	*service.ServiceContext
}

// RegistrationUser godoc
// @Summary Register a new user
// @Description Register a new user with the input payload
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body models.User true "Register user"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/register [post]
func (l *template) RegistrationUser(c echo.Context) error {
	var reqUser models.User
	if err := c.Bind(&reqUser); err != nil {
		res := utils.Response{
			Code: http.StatusInternalServerError,
			Err:  utils.STATUS_INTERNAL_ERR,
			Msg:  err.Error(),
		}
		return c.JSON(res.Code, res)
	}

	utils.PrintStruct(reqUser)

	if err := reqUser.Validate(); err != nil {
		res := utils.Response{
			Code: http.StatusBadRequest,
			Err:  utils.STATUS_BAD_REQUEST,
			Msg:  err.Error(),
		}
		return c.JSON(res.Code, res)
	}

	if err := l.TemplateService.InsertUser(c.Request().Context(), reqUser); err != nil {
		res := utils.Response{
			Code: http.StatusInternalServerError,
			Err:  utils.STATUS_INTERNAL_ERR,
			Msg:  err.Error(),
		}
		return c.JSON(res.Code, res)
	}

	return c.JSON(http.StatusOK, utils.Response{
		Msg: "Success Registration.",
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate a user and return a token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param credentials body models.User true "User credentials"
// @Success 200 {object} utils.Response{data=models.LoginResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/login [post]
func (l *template) Login(c echo.Context) error {
	var req models.User
	if err := c.Bind(&req); err != nil {
		res := utils.Response{
			Code: http.StatusInternalServerError,
			Err:  utils.STATUS_INTERNAL_ERR,
			Msg:  err.Error(),
		}
		return c.JSON(res.Code, res)
	}

	if err := req.ValidateLogin(); err != nil {
		res := utils.Response{
			Code: http.StatusBadRequest,
			Err:  utils.STATUS_BAD_REQUEST,
			Msg:  err.Error(),
		}
		return c.JSON(res.Code, res)
	}

	userData, err := l.TemplateService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		res := utils.Response{
			Code: http.StatusInternalServerError,
			Err:  utils.STATUS_INTERNAL_ERR,
			Msg:  err.Error(),
		}
		return c.JSON(res.Code, res)
	}

	token, err := l.JWTService.GenerateJWTToken(c, models.JWTRequest{
		ID:    int(userData.ID),
		Email: userData.Email,
		Name:  userData.Username,
	})

	if err != nil {
		res := utils.Response{
			Code: http.StatusInternalServerError,
			Err:  utils.STATUS_INTERNAL_ERR,
			Msg:  err.Error(),
		}
		return c.JSON(res.Code, res)
	}

	return c.JSON(http.StatusOK, utils.Response{
		Data: models.LoginResponse{
			Id:    int(userData.ID),
			Token: token,
		},
		Msg: "Success Login.",
	})
}

// Profile godoc
// @Summary Get user profile
// @Description Get the profile of the authenticated user
// @Tags auth
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} utils.Response{data=models.UserResponse}
// @Failure 500 {object} utils.Response
// @Router /auth/profile [get]
func (l *template) Profile(c echo.Context) error {
	token := middleware.GetTokenFromContext(c)
	userData, err := l.TemplateService.GetUserByID(c.Request().Context(), token.ID)
	if err != nil {
		res := utils.Response{
			Code: http.StatusInternalServerError,
			Err:  utils.STATUS_INTERNAL_ERR,
			Msg:  err.Error(),
		}
		return c.JSON(res.Code, res)
	}

	return c.JSON(http.StatusOK, utils.Response{
		Data: userData,
		Msg:  "Success Get Profile.",
	})
}
