package middleware

import (
	"log"
	"net/http"
	"os"
	"sid/service/handler"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Upgrade(validator func(*gin.Context, *jwt.RegisteredClaims) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, e := upgrader.Upgrade(c.Writer, c.Request, nil)
		if e != nil {
			log.Printf("Upgrade error: %+v\n", e)
			c.Abort()
			return
		}
		defer conn.Close()

		e = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if e != nil {
			log.Printf("SetReadDeadline error: %+v\n", e)
			c.Abort()
			return
		}

		msg := &struct {
			Authorization string `json:"authorization"`
		}{}
		e = conn.ReadJSON(msg)
		if e != nil {
			log.Printf("ReadJSON error: %+v\n", e)
			e = conn.WriteJSON(handler.ErrUnauthorized)
			if e != nil {
				log.Printf("WriteJSON error: %+v\n", e)
			}
			c.Abort()
			return
		}

		// token := c.GetHeader("Authorization")
		token := msg.Authorization

		claim := &jwt.RegisteredClaims{}
		info, e := jwt.ParseWithClaims(token, claim, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if e != nil || !info.Valid {
			e = conn.WriteJSON(handler.ErrUnauthorized)
			if e != nil {
				log.Printf("WriteJSON error: %+v\n", e)
			}
			c.Abort()
			return
		}

		e = validator(c, claim)

		if e != nil || !info.Valid {
			if e != nil {
				log.Printf("WriteJSON error: %+v\n", e)
			}
			c.Abort()
			return
		}

		c.Set("Websocket", conn)

		c.Next()
	}
}
