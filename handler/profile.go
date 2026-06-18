package handler

import (
	"net/http"
	"sid/service/feature"

	"github.com/gin-gonic/gin"
)

type Profile struct {
	feature *feature.Profile
}

func NewProfile(feature *feature.Profile) *Profile {
	return &Profile{feature: feature}
}

func (u *Profile) ProfileSignIn(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"method": "GET"})
}

func (u *Profile) ProfileSignUp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"method": "POST"})
}

func (u *Profile) ProfileSignOut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"method": "PUT"})
}

func (u *Profile) ProfileAuth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"method": "DELETE"})
}
