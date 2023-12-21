package server

import (
	docs "github.com/gweebg/probum-users/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/gweebg/probum-users/controllers"
)

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

			userGroup.POST("/", user.CreateUser)
			userGroup.GET("/:id", user.GetUser)
			userGroup.PATCH("/:id", user.UpdateUser)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}
