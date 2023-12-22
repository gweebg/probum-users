package controllers

import (
	"github.com/gweebg/probum-users/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gweebg/probum-users/forms"
	"github.com/gweebg/probum-users/models"
)

type UserController struct{}

var userModel = new(models.User)

// GetUser       godoc
// @Summary      Retrieve a user from the database.
// @Description  Retrieves the user with the specified school id, as a json object.
// @Tags         users
// @Produce      json
// @Param        id  path      string  true  "search user by id"
// @Success      200   {object}  models.User
// @Router       /user/{id} [get]
func (u UserController) GetUser(c *gin.Context) {

	userId := c.Param("id")
	if userId != "" {

		user, err := userModel.Get(userId)
		if err != nil {
			utils.HandleAbort(c, http.StatusNotFound, "user not found", err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
		return
	}

	utils.HandleAbort(c, http.StatusInternalServerError, "could not add user to the database", "")
	return
}

// CreateUser    godoc
// @Summary      Insert a user into the system.
// @Description  When provided with a user object, this endpoint inserts it into the database if the id does not exist and the object is well formed.
// @Tags         users
// @Produce      json
// @Param        user  body      forms.UserSignup  true  "signup user form"
// @Success      200   {object}  models.User
// @Router       /user [post]
func (u UserController) CreateUser(c *gin.Context) {

	// Decode form body.
	var newUser forms.UserSignup
	err := c.ShouldBindJSON(&newUser)

	if err != nil {
		utils.HandleAbort(c, http.StatusBadRequest, "could not unmarshal request body into 'forms.UserSignup'", err.Error())
		return
	}

	// Send sign up request to authentication service.
	cookie, _ := c.Cookie("Authorization") // cookie is set (verified at the middleware)

	headers := map[string]string{
		"Authorization": cookie,
	}
	payload := struct {
		UId      string
		Password string
	}{
		UId:      newUser.UId,
		Password: newUser.Password,
	}

	_, err = utils.SendHTTPRequest(
		c.GetString("endpoints.auth.signup.method"),
		c.GetString("endpoints.auth.base")+c.GetString("endpoints.auth.signup.uri"),
		headers, payload,
	)
	if err != nil {
		utils.HandleAbort(c, http.StatusInternalServerError, "could not reach authentication server", err.Error())
		return
	}

	// Insert the new user into the database.
	user, err := userModel.Create(newUser)
	if err != nil {
		utils.HandleAbort(c, http.StatusInternalServerError, "could not add user to the database", err.Error())
		return
	}

	// Complete.
	c.JSON(http.StatusCreated, gin.H{
		"user": user,
	})
}

func (u UserController) UpdateUser(c *gin.Context) {

	userId := c.Param("id")
	if userId != "" {

		var userUpdate forms.UserUpdate
		err := c.ShouldBindJSON(&userUpdate)

		if err != nil {
			utils.HandleAbort(c, http.StatusBadRequest, "could not unmarshal request body into 'forms.UserUpdate'", err.Error())
			return
		}

		user, err := userModel.Update(userId, userUpdate)
		if err != nil {
			utils.HandleAbort(c, http.StatusInternalServerError, "could not update user information", err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
		return
	}

	// todo: add communication with the auth microservice

	utils.HandleAbort(c, http.StatusBadRequest, "field 'id' is not specified", "")
	return
}
