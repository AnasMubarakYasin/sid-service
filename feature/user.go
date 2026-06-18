package feature

import (
	"context"
	"sid/service/model"
	"sid/service/repository"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	u *repository.User
	p *repository.Profile
}

func NewUser(u *repository.User, p *repository.Profile) *User {
	user := User{u, p}
	return &user
}

type ParamSignUp struct {
	Username string         `form:"username" json:"username" xml:"username" binding:"required"`
	Role     model.UserRole `form:"role" json:"role" xml:"role" binding:"required"`
	Password string         `form:"password" json:"password" xml:"password" binding:"required"`
	Confirm  string         `form:"confirm" json:"confirm" xml:"confirm" binding:"required"`
}

type ParamSignIn struct {
	Username string `form:"username" json:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

var (
	ErrConfirmMismatch = &Error{Code: "Confirm Mismatch", Message: "Password and Confirmation mismatch"}
	ErrUnameExists     = &Error{Code: "Uname Exists", Message: "Username already exists"}
	ErrUnameNoExists   = &Error{Code: "Uname No Exists", Message: "Username not exists"}
	ErrPwdMismatch     = &Error{Code: "Password Mismatch", Message: "Password mismatch"}
)

func (f *User) UserSignUp(c context.Context, p *ParamSignUp) (r *model.User, e error) {
	d, e := f.u.UserFindByName(c, p.Username)

	if e != nil {
		return nil, e
	}

	if d != nil {
		return nil, ErrUnameExists
	}

	if p.Password != p.Confirm {
		return nil, ErrConfirmMismatch
	}

	pwd, e := bcrypt.GenerateFromPassword(([]byte(p.Password)), bcrypt.DefaultCost)

	if e != nil {
		return nil, e
	}

	r, e = f.u.UserCreate(c, &model.User{Name: p.Username, Role: p.Role, Password: string(pwd)})

	if e != nil {
		return nil, e
	}

	r.Password = ""

	return r, nil
}

func (f *User) UserSignIn(c context.Context, p *ParamSignIn) (r *model.User, e error) {
	r, e = f.u.UserFindByName(c, p.Username)

	if e != nil {
		return nil, e
	}

	if r == nil {
		return nil, ErrUnameNoExists
	}

	hpwd := []byte(r.Password)
	pwd := []byte(p.Password)

	if e = bcrypt.CompareHashAndPassword(hpwd, pwd); e != nil {
		return nil, ErrPwdMismatch
	}

	r.Password = ""

	return r, nil
}

func (f *User) UserSignOut(c context.Context) {
}

func (f *User) UserAuth(c context.Context, n string) (r *model.User, e error) {
	r, e = f.u.UserFindByName(c, n)

	if e != nil {
		return nil, e
	}

	if r == nil {
		return nil, ErrUnameNoExists
	}

	r.Password = ""

	return r, nil
}
