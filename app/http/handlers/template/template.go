package template

import (
	"base-api/data/models"
	"base-api/infra/context/service"
	"base-api/infra/middleware"
	"base-api/utils"
	"encoding/json"
	"net/http"
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
func (l *template) RegistrationUser(w http.ResponseWriter, r *http.Request) {
	var reqUser models.User
	err := json.NewDecoder(r.Body).Decode(&reqUser)
	if err != nil {
		res := utils.Response{
			Code: http.StatusInternalServerError,
			Data: nil,
			Err:  utils.STATUS_INTERNAL_ERR,
			Msg:  err.Error(),
		}
		res.JSONResponse(w)
		return
	}

	utils.PrintStruct(reqUser)

	err = reqUser.Validate()
	if err != nil {
		res := utils.Response{
			Code: http.StatusBadRequest,
			Data: nil,
			Err:  utils.STATUS_BAD_REQUEST,
			Msg:  err.Error(),
		}
		res.JSONResponse(w)
		return
	}

	err = l.TemplateService.InsertUser(r.Context(), reqUser)
	if err != nil {
		res := utils.Response{
			Code: http.StatusInternalServerError,
			Data: nil,
			Err:  utils.STATUS_INTERNAL_ERR,
			Msg:  err.Error(),
		}
		res.JSONResponse(w)
		return
	}

	res := utils.Response{
		Code: http.StatusOK,
		Msg:  "Success Registration.",
	}
	res.JSONResponse(w)
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
func (l *template) Login(w http.ResponseWriter, r *http.Request) {
	var req models.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		res := utils.Response{
			Code: http.StatusInternalServerError,
			Data: nil,
			Err:  utils.STATUS_INTERNAL_ERR,
			Msg:  err.Error(),
		}
		res.JSONResponse(w)
		return
	}

	err = req.ValidateLogin()
	if err != nil {
		res := utils.Response{
			Code: http.StatusBadRequest,
			Data: nil,
			Err:  utils.STATUS_BAD_REQUEST,
			Msg:  err.Error(),
		}
		res.JSONResponse(w)
		return
	}

	userData, err := l.TemplateService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		res := utils.Response{
			Code: http.StatusInternalServerError,
			Data: nil,
			Err:  utils.STATUS_INTERNAL_ERR,
			Msg:  err.Error(),
		}
		res.JSONResponse(w)
		return
	}

	token, err := l.JWTService.GenerateJWTToken(r.Context(), models.JWTRequest{
		ID:    int(userData.ID),
		Email: userData.Email,
		Name:  userData.Username,
	})

	if err != nil {
		res := utils.Response{
			Code: http.StatusInternalServerError,
			Data: nil,
			Err:  utils.STATUS_INTERNAL_ERR,
			Msg:  err.Error(),
		}
		res.JSONResponse(w)
		return
	}

	res := utils.Response{
		Code: http.StatusOK,
		Data: models.LoginResponse{
			Id:    int(userData.ID),
			Token: token,
		},
		Msg: "Success Login.",
	}
	res.JSONResponse(w)
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
func (l *template) Profile(w http.ResponseWriter, r *http.Request) {
	token := middleware.GetTokenFromContext(r.Context())
	userData, err := l.TemplateService.GetUserByID(r.Context(), token.ID)
	if err != nil {
		res := utils.Response{
			Code: http.StatusInternalServerError,
			Data: nil,
			Err:  utils.STATUS_INTERNAL_ERR,
			Msg:  err.Error(),
		}
		res.JSONResponse(w)
		return
	}

	res := utils.Response{
		Code: http.StatusOK,
		Data: userData,
		Msg:  "Success Get Profile.",
	}
	res.JSONResponse(w)
}
