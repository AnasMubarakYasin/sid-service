package middleware

import (
	"errors"
	"net/http"
	"sid/service/handler"

	"github.com/gin-gonic/gin"
)

// Error captures errors and returns a consistent JSON error response
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		var hterr *handler.HTERR
		if errors.As(err, &hterr) {
			c.JSON(hterr.Status, gin.H{
				"success": false,
				"error":   gin.H{"code": hterr.Code, "message": hterr.Message},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   gin.H{"code": "INTERNAL", "message": err.Error()},
			})
		}
	}
}
