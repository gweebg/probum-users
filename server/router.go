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

			userGroup.POST("/login", user.Login)

			userGroup.Use(middlewares.RequireAuth)
			{
				userGroup.GET("/:id", user.Get)
				userGroup.GET("/", user.Get)
				userGroup.PATCH("/", user.Update)
				userGroup.POST("/", user.Create) // Needs to be protected, only Admins | Techs can create.
			}

		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return router
}
