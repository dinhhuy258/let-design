package middleware

import (
	"errors"
	"job-service/config"
	"job-service/internal/entity"
	"job-service/internal/usecase"
	"job-service/pkg/logger"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var (
	identityKey = jwt.IdentityKey
	userIDKey   = "user_id"
	userData    = "user_data"
)

func NewAuthMiddleware(
	authUsecase usecase.AuthUsecase,
	conf *config.Config,
	logger *logger.Logger,
) (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "job-service",
		Key:         []byte(conf.Jwt.SecretKey),
		Timeout:     conf.Jwt.AccessTokenTimeOut,
		MaxRefresh:  conf.Jwt.RefreshTokenTimeOut,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if user, ok := data.(*entity.User); ok {
				return jwt.MapClaims{
					userIDKey: user.Id,
				}
			}

			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)

			return &entity.User{
				Id: claims[userIDKey].(uint64),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var login entity.Login
			if err := c.ShouldBind(&login); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			user, err := authUsecase.AttemptLogin(c, login.Username, login.Password)
			if err != nil {
				if errors.Is(err, entity.ErrAuthFailed) || errors.Is(err, entity.ErrUserNotFound) {
					return "", jwt.ErrFailedAuthentication
				}

				return "", err
			}

			c.Set(userData, user)

			return user, nil
		},
		Authorizator: func(data interface{}, _ *gin.Context) bool {
			if _, ok := data.(*entity.User); ok {
				return true
			}

			return false
		},
		LoginResponse: func(c *gin.Context, _ int, token string, expire time.Time) {
			u, _ := c.Get(userData)
			user, _ := u.(*entity.User)

			c.JSON(http.StatusOK, gin.H{
				"expire": expire,
				"token":  token,
				"user":   user,
			})
		},
		RefreshResponse: func(c *gin.Context, _ int, token string, expire time.Time) {
			u, _ := c.Get(userData)
			user, _ := u.(*entity.User)

			c.JSON(http.StatusOK, gin.H{
				"expire": expire,
				"token":  token,
				"user":   user,
			})
		},
		LogoutResponse: func(c *gin.Context, code int) {
			c.Status(code)
		},
		Unauthorized: func(c *gin.Context, _ int, message string) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		logger.Error("Failed to init jwt middleware %v", err)

		return nil, err
	}

	return authMiddleware, nil
}
