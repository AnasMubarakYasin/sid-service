package feature

import (
	"sid/service/model"
	"sid/service/repository"

	"github.com/gin-gonic/gin"
)

type Profile struct {
	p *repository.Profile
}

func NewProfile(p *repository.Profile) *Profile {
	Profile := Profile{p}
	return &Profile
}

func (f *Profile) ProfileSignIn(c *gin.Context) {
}

func (f *Profile) ProfileSignUp(c *gin.Context, p *ParamSignUp) (*model.Profile, error) {
	return nil, nil
}

func (f *Profile) ProfileSignOut(c *gin.Context) {
}

func (f *Profile) ProfileAuth(c *gin.Context) {
}
