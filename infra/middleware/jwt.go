package middleware

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"

	"base-api/config"
	"base-api/constants"
	"base-api/data/models"

	"github.com/dgrijalva/jwt-go"
)

const (
	bearer      = "Bearer"
	AuthKeyUser = "auth-user"
)

var (
	JWTSigningMethod = jwt.SigningMethodHS256
)

type JWTClaims struct {
	jwt.StandardClaims
	ID           int
	Email        string
	Name         string
	CostCenterID int
	JoinDate     *time.Time
	Role         string
}

type JWTInterface interface {
	ExtractJWTClaims(c echo.Context) (claims *JWTClaims, err error)
	ValidateTokenIssuer(claims *JWTClaims) error
	ValidateTokenExpire(c echo.Context, claims *JWTClaims) error
	GenerateJWTToken(c echo.Context, request models.JWTRequest) (string, error)
}

type jwtObj struct {
	config *config.JWTConfig
	redis  *redis.Client
}

func NewJWT(cfg *config.JWTConfig, redis *redis.Client) JWTInterface {
	return &jwtObj{
		config: cfg,
		redis:  redis,
	}
}

func (j *jwtObj) ExtractJWTClaims(c echo.Context) (claims *JWTClaims, err error) {
	tokenHeader := c.Request().Header.Get(constants.Authorization)
	if tokenHeader == "" {
		return nil, constants.ErrTokenIsRequired
	}

	splitToken := strings.Split(tokenHeader, bearer)
	if len(splitToken) != 2 {
		return nil, constants.ErrTokenIsRequired
	}
	reqToken := strings.TrimSpace(splitToken[1])

	t, err := jwt.ParseWithClaims(reqToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.config.Secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims = t.Claims.(*JWTClaims)

	if err := j.ValidateTokenIssuer(claims); err != nil {
		return nil, err
	}

	if err := j.ValidateTokenExpire(c, claims); err != nil {
		return nil, err
	}

	return claims, nil
}

// ValidateTokenIssuer is for validate token issuer
func (j *jwtObj) ValidateTokenIssuer(claims *JWTClaims) (err error) {
	if claims.Issuer != j.config.Issuer {
		return constants.ErrTokenInvalid
	}
	return nil
}

// ValidateTokenExpire is for validate Token Expire
func (j *jwtObj) ValidateTokenExpire(c echo.Context, claims *JWTClaims) error {
	key := AuthKeyUser
	token, err := j.GetTokenFromRedis(c.Request().Context(), claims.ID, key)
	if err != nil {
		log.Printf("%s: %s", constants.ErrGetTokenFromRedis, err.Error())
		return constants.ErrTokenAlreadyExpired
	}

	reqToken := strings.TrimPrefix(c.Request().Header.Get(constants.Authorization), "Bearer ")
	if token != reqToken {
		return constants.ErrTokenReplaced
	}

	return nil
}

func (j *jwtObj) GetTokenFromRedis(ctx context.Context, id int, authKey string) (string, error) {
	key := fmt.Sprintf("%s:%d", authKey, id)
	val, err := j.redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func (j *jwtObj) DeleteTokenFromRedis(ctx context.Context, id int, authKey string) error {
	key := fmt.Sprintf("%s:%d", authKey, id)
	_, err := j.redis.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}

func (j *jwtObj) GenerateJWTToken(ctx echo.Context, request models.JWTRequest) (string, error) {
	JWTSignatureKey := []byte(j.config.Secret)
	expireTime := time.Now().Add(time.Duration(j.config.TokenLifeTimeHour) * time.Hour)

	key := AuthKeyUser
	claims := JWTClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    j.config.Issuer,
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		ID:    request.ID,
		Email: request.Email,
		Name:  request.Name,
	}

	token := jwt.NewWithClaims(
		JWTSigningMethod,
		claims,
	)

	// create token client
	signedToken, err := token.SignedString(JWTSignatureKey)
	if err != nil {
		return "", err
	}

	// Save token to redis
	err = j.SaveTokenToRedis(ctx.Request().Context(), request.ID, j.config.TokenLifeTimeHour, signedToken, key)
	if err != nil {
		err = constants.ErrSaveTokenToRedis
		return "", err
	}
	return signedToken, nil
}

func (j *jwtObj) SaveTokenToRedis(ctx context.Context, id, hour int, token, authKey string) error {
	key := fmt.Sprintf("%s:%d", authKey, id)
	ttl := time.Duration(hour) * time.Hour
	err := j.redis.Set(ctx, key, token, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}
