package server

import (
	docs "github.com/gweebg/probum-users/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/gweebg/probum-users/controllers"
	"github.com/gweebg/probum-users/middlewares"
)

// @title           Probum User Manager
// @version         1.0
// @description     A user management service for the application Probum.

// @contact.name   Guilherme
// @contact.url    https://github.com/gweebg

// @license.name  MIT
// @license.url   https://mit-license.org/

// @host      localhost:3000
// @BasePath  /api/v1

func NewRouter() *gin.Engine {

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	docs.SwaggerInfo.BasePath = "/api/v1"

	v1 := router.Group("api/v1")
	{
		userGroup := v1.Group("user")
		{
			user := new(controllers.UserController)

			userGroup.GET("/:id", user.Get)

			userGroup.Use(middlewares.RequireAuth)
			{
				userGroup.GET("/", user.GetCurrent)
				userGroup.POST("/", user.Create)
				userGroup.PATCH("/", user.Update)
			}
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return router
}
