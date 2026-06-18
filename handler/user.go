package handler

import (
	"os"
	"sid/service/feature"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	f *feature.User
}

func NewUser(feature *feature.User) *User {
	return &User{f: feature}
}

func (h *User) UserSignUp(c *gin.Context) {
	p := &feature.ParamSignUp{}

	if e := c.ShouldBindJSON(p); e != nil {
		_ = c.Error(ErrBadRequest)
		return
	}

	r, e := h.f.UserSignUp(c, p)

	if e != nil {
		Fail(c, 500, e)
		return
	}

	t, e := GenToken(r.Name)

	if e != nil {
		_ = c.Error(ErrInternal)
		return
	}

	OKM(c, r, &Meta{Token: t})
}

func (h *User) UserSignIn(c *gin.Context) {
	p := &feature.ParamSignIn{}

	if e := c.ShouldBindJSON(p); e != nil {
		_ = c.Error(ErrBadRequest)
		return
	}

	r, e := h.f.UserSignIn(c, p)

	if e != nil {
		Fail(c, 500, e)
		return
	}

	t, e := GenToken(r.Name)

	if e != nil {
		_ = c.Error(ErrInternal)
		return
	}

	OKM(c, r, &Meta{Token: t})
}

func (h *User) UserSignOut(c *gin.Context) {
	OK(c, &Response{})
}

func (h *User) UserAuth(c *gin.Context) {
	header := new(HAuth)
	if err := c.ShouldBindHeader(header); err != nil {
		_ = c.Error(ErrUnauthorized)
		return
	}
	// token := c.GetHeader("Authorization")
	token := header.Authorization

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	cl := &jwt.RegisteredClaims{}
	v, e := VerifyToken(token, cl)

	if e != nil || !v {
		_ = c.Error(ErrUnauthorized)
		return
	}

	r, e := h.f.UserAuth(c, cl.Subject)

	if e != nil {
		_ = c.Error(ErrInternal)
		return
	}

	OK(c, r)
}

func (h *User) UserAuthz(c *gin.Context) {

	c.Param("p")
	r, o := c.Get("User")

	if !o {
		_ = c.Error(ErrInternal)
		return
	}

	OK(c, r)
}

func (h *User) UserValid(c *gin.Context, cl *jwt.RegisteredClaims) error {
	r, e := h.f.UserAuth(c, cl.Subject)

	if e != nil {
		return e
	}

	c.Set("User", r)

	return nil
}

type HAuth struct {
	Authorization string `header:"Authorization" binding:"required"`
}

func GenToken(subject string) (string, error) {
	// Create token
	expr := time.Now().Add(24 * time.Hour)
	claims := &jwt.RegisteredClaims{
		Subject:   subject,
		ExpiresAt: jwt.NewNumericDate(expr),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	return signed, err
}

func VerifyToken(token string, claims *jwt.RegisteredClaims) (bool, error) {
	i, e := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if e != nil {
		return false, e
	}
	return i.Valid, e
}
