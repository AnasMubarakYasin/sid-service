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

	api := router.Group("api")
	event := router.Group("event")
	message := router.Group("message")

	router.GET("/", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "hello"}) })

	api.POST("/user/sign-up", user.UserSignUp)
	api.POST("/user/sign-in", user.UserSignIn)
	api.POST("/user/sign-out", user.UserSignOut)
	api.GET("/user/auth", user.UserAuth)
	api.GET("/user/authz/*page", middleware.Auth(user.UserValid), user.UserAuthz)

	{
		group := api.Group("experiment")
		group.Use(middleware.Auth(user.UserValid))
		group.POST("/", experiment.Create)
		group.PUT("/:id", experiment.Update)
		group.DELETE("/:id", experiment.Delete)
		group.GET("/:id", experiment.Get)
		group.GET("/", experiment.List)

		group.POST("/lab", experiment.CreateLab)
		group.DELETE("/lab/:id", experiment.RemoveLab)
		event.GET("/experiment", middleware.Auth(user.UserValid), experiment.Event)
		message.GET("/experiment/:id", middleware.Upgrade(user.UserValid), experiment.Collaboration)
	}

	{
		group := api.Group("subject")
		group.Use(middleware.Auth(user.UserValid))
		group.POST("/", subject.Create)
		group.PUT("/:id", subject.Update)
		group.DELETE("/:id", subject.Delete)
		group.GET("/:id", subject.Get)
		group.GET("/", subject.List)
	}

}
