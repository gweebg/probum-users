package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gweebg/probum-users/forms"
	"github.com/gweebg/probum-users/models"
)

type UserController struct{}

var userModel = new(models.User)

func (u UserController) GetUser(c *gin.Context) {

	userId := c.Param("id")
	if userId != "" {

		user, err := userModel.Get(userId)
		if err != nil {

			c.JSON(http.StatusNotFound, gin.H{
				"message": "user not found",
				"error":   err.Error(),
			})
			c.Abort()
			return

		}

		// todo: add communication with authentication server

		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "bad request, 'id' not specified",
	})
	c.Abort()
	return

}

func (u UserController) CreateUser(c *gin.Context) {

	var newUser forms.UserSignup
	err := c.ShouldBindJSON(&newUser)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "could not unmarshal request body into 'forms.UserSignup'",
			"error":   err.Error(),
		})
		c.Abort()
		return
	}

	user, err := userModel.Create(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not add user to the database",
			"error":   err.Error(),
		})
		c.Abort()
		return
	}

	// todo: add communication with the auth microservice

	c.JSON(http.StatusCreated, gin.H{
		"user": user,
	})
	return
}

func (u UserController) UpdateUser(c *gin.Context) {

	userId := c.Param("id")
	if userId != "" {

		var userUpdate forms.UserUpdate
		err := c.ShouldBindJSON(&userUpdate)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "could not unmarshal request body into 'forms.UserUpdate'",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		user, err := userModel.Update(userId, userUpdate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "could not update user information",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
		return
	}

	// todo: add communication with the auth microservice

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "'id' is not specified",
	})
	c.Abort()
	return
}
