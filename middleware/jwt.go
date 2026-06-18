package middleware

import (
	"os"
	"sid/service/handler"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(validator func(*gin.Context, *jwt.RegisteredClaims) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := &struct {
			Authorization string `header:"Authorization" binding:"required"`
		}{}
		if err := c.ShouldBindHeader(header); err != nil {
			_ = c.AbortWithError(handler.ErrUnauthorized.Status, handler.ErrUnauthorized)
			return
		}

		// token := c.GetHeader("Authorization")
		token := header.Authorization
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		claim := &jwt.RegisteredClaims{}
		info, err := jwt.ParseWithClaims(token, claim, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil || !info.Valid {
			_ = c.AbortWithError(handler.ErrUnauthorized.Status, handler.ErrUnauthorized)
			return
		}

		err = validator(c, claim)

		if err != nil || !info.Valid {
			_ = c.AbortWithError(handler.ErrUnauthorized.Status, handler.ErrUnauthorized)
			return
		}
		c.Next()
	}
}
