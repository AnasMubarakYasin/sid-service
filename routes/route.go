package routes

import (
	"log/slog"
	"sid/service/feature"
	"sid/service/handler"
	"sid/service/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine, logger *slog.Logger, features *feature.Core) {
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	// router.Use(middleware.SlogLogger(logger))
	router.Use(middleware.RateLimiter())
	router.Use(middleware.CORS())
	router.Use(middleware.Header())
	router.Use(middleware.ErrorHandler())

	user := handler.NewUser(features.User)
	subject := handler.NewSubject(features.Subject)
	experiment := handler.NewExperiment(features.Experiment)

	router.GET("/", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "hello"}) })

	router.POST("/user/sign-up", user.UserSignUp)
	router.POST("/user/sign-in", user.UserSignIn)
	router.POST("/user/sign-out", user.UserSignOut)
	router.GET("/user/auth", user.UserAuth)
	router.GET("/user/authz/*page", middleware.Auth(user.UserValid), user.UserAuthz)

	{
		group := router.Group("experiment")
		group.Use(middleware.Auth(user.UserValid))
		group.POST("/", experiment.Create)
		group.PUT("/:id", experiment.Update)
		group.DELETE("/:id", experiment.Delete)
		group.GET("/:id", experiment.Get)
		group.GET("/", experiment.List)

		group.POST("/lab", experiment.CreateLab)
		group.DELETE("/lab/:id", experiment.RemoveLab)
		router.GET("/event/experiment", middleware.Auth(user.UserValid), experiment.Event)
		router.GET("/message/experiment/:id", middleware.Upgrade(user.UserValid), experiment.Collaboration)
	}

	{
		group := router.Group("subject")
		group.Use(middleware.Auth(user.UserValid))
		group.POST("/", subject.Create)
		group.PUT("/:id", subject.Update)
		group.DELETE("/:id", subject.Delete)
		group.GET("/:id", subject.Get)
		group.GET("/", subject.List)
	}

}
