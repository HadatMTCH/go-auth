package middleware

import (
	"base-api/constants"
	"base-api/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type TokenMiddlewareInterface interface {
	TokenAuthorize() echo.MiddlewareFunc
}

type tokenMiddlewareObj struct {
	JWTService JWTInterface
}

func NewTokenMiddleware(jwt JWTInterface) TokenMiddlewareInterface {
	return &tokenMiddlewareObj{jwt}
}

func (m *tokenMiddlewareObj) TokenAuthorize() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := m.JWTService.ExtractJWTClaims(c)
			if err != nil {
				code := http.StatusInternalServerError
				errMsg := utils.STATUS_INTERNAL_ERR

				if err == constants.ErrTokenAlreadyExpired ||
					err == constants.ErrTokenReplaced ||
					err == constants.ErrTokenInvalid {
					code = http.StatusUnauthorized
					errMsg = utils.STATUS_UNAUTHORIZED
				}

				log.Error(err)
				return c.JSON(code, utils.Response{
					Code: code,
					Err:  errMsg,
					Msg:  err.Error(),
				})
			}

			// Store claims in Echo context
			c.Set("user", claims)
			return next(c)
		}
	}
}

// GetTokenFromContext now uses Echo's context
func GetTokenFromContext(c echo.Context) *JWTClaims {
	claims, ok := c.Get("user").(*JWTClaims)
	if !ok {
		return nil
	}
	return claims
}
