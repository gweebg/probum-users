package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gweebg/probum-users/forms"
	"github.com/gweebg/probum-users/models"
)

// @BasePath /api/v1

type UserController struct{}

var userModel = new(models.User)

// GetUser retrieves a user from the database using their school identification.
// @Summary Retrieve User by School ID
// @Description Retrieves a user from the database based on the school identification (pg|a[0-9]\d{10}).
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Router /user/{id} [get]
func (u UserController) GetUser(c *gin.Context) {

	userId := c.Param("id")
	if userId != "" {

		user, err := userModel.GetUserById(userId)
		if err != nil {

			c.JSON(http.StatusNotFound, gin.H{
				"message": "user not found",
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

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "bad request, 'id' not specified",
	})
	c.Abort()
	return

}

func (u UserController) CreateUser(c *gin.Context) {
	return
}

func updateFormHandler(c *gin.Context) (forms.UserUpdate, error) {

	var info forms.UserUpdate
	err := c.ShouldBindJSON(&info)

	return info, err
}

func (u UserController) UpdateUser(c *gin.Context) {

	userId := c.Param("id")
	if userId != "" {

		updateForm, err := updateFormHandler(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "could not unmarshal json payload into 'forms.UserUpdate'",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		user, err := userModel.Update(userId, updateForm)
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

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "'id' is not specified",
	})
	c.Abort()
	return

}
