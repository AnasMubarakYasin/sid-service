package handler

import (
	"errors"
	"net/http"
	"sid/service/feature"
	"sid/service/model"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Response is the standard API envelope.
type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   *Error `json:"error,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Page       int    `json:"page,omitempty"`
	PerPage    int    `json:"per_page,omitempty"`
	Total      int    `json:"total,omitempty"`
	TotalPages int    `json:"total_pages,omitempty"`
	Token      string `json:"token,omitempty"`
}

type ParamID struct {
	ID int64 `uri:"id" binding:"required"`
}

func GetUser(c *gin.Context) (*model.User, error) {
	cusr, ok := c.Get("User")
	if !ok {
		return nil, ErrNotFound
	}
	usr, ok := cusr.(*model.User)
	if !ok {
		return nil, ErrNotFound
	}
	return usr, nil
}

func GetWebsocket(c *gin.Context) (*websocket.Conn, error) {
	ws, ok := c.Get("Websocket")
	if !ok {
		return nil, ErrInternal
	}
	cast, ok := ws.(*websocket.Conn)
	if !ok {
		return nil, ErrInternal
	}
	return cast, nil
}

// OK sends a success response.
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &Response{
		Success: true,
		Data:    data,
	})
}

func OKM(c *gin.Context, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, &Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// Fail sends an error response.
func Fail(c *gin.Context, s int, e error) {
	var fe *feature.Error
	if errors.As(e, &fe) {
		c.JSON(s, gin.H{
			"success": false,
			"error":   gin.H{"code": fe.Code, "message": fe.Message},
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   gin.H{"code": "INTERNAL", "message": e.Error()},
		})
	}
}

// HTERR represents a structured API error.
type HTERR struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *HTERR) Error() string {
	return e.Message
}

var (
	ErrNotFound     = &HTERR{Status: 404, Code: "NOT_FOUND", Message: "resource not found"}
	ErrUnauthorized = &HTERR{Status: 401, Code: "UNAUTHORIZED", Message: "authentication required"}
	ErrBadRequest   = &HTERR{Status: 400, Code: "BAD_REQUEST", Message: "invalid request"}

	ErrInternal = &HTERR{Status: 500, Code: "INTERNAL_SERVER_ERROR", Message: "server error"}
	// ErrUnAuth = &HTERR{Status: 511, Code: "StatusNetworkAuthenticationRequired", Message: "server error"}
)
