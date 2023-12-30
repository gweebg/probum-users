package controllers

import (
	"fmt"
	"github.com/gweebg/probum-users/config"
	"github.com/gweebg/probum-users/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gweebg/probum-users/forms"
	"github.com/gweebg/probum-users/models"
)

type UserController struct{}

var userModel = new(models.User)

// Get           godoc
// @Summary      Retrieve a user from the database
// @Description  Retrieves the user (as a json object) with the specified school id. A school id follows the expression (pg|a)[1-9]\d{6}.
// @Tags         users
// @Produce      json
// @Param        id  path      string  true  "User id"
// @Success      200   {object}  models.User
// @Router       /user/{id} [get]
func (u UserController) Get(c *gin.Context) {

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

// GetCurrent    godoc
// @Summary      Retrieve the current authenticated user.
// @Description  Retrieves the user corresponding to the provided authentication jwt token.
// @Tags         users
// @Produce      json
// @Success      200   {object}  models.User
// @Security 	 BearerToken
// @Router       /user [get]
func (u UserController) GetCurrent(c *gin.Context) {

	// If it got to this point, then 'user' is set in *gin.Context.
	user, exists := c.Get("user")
	if !exists {
		utils.HandleAbort(c, http.StatusNotFound, "user not found, even though is authenticated", "")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user.(*models.User),
	})
	return
}

// Create        godoc
// @Summary      Insert a user into the system
// @Description  When provided with a user object, this endpoint inserts it into the database if the id does not exist and the object is well-formed.
// @Tags         users
// @Produce      json
// @Param        user  body      forms.UserSignup  true  "User signup json object"
// @Success      200   {object}  models.User
// @Security 	 BearerToken
// @Router       /user [post]
func (u UserController) Create(c *gin.Context) {

	requester, _ := c.Get("user")
	reqUser := requester.(*models.User)

	// Check roles.
	if reqUser.Role != "tech" && reqUser.Role != "admin" {
		utils.HandleAbort(c, http.StatusUnauthorized, fmt.Sprintf("role '%s' is not allowed to create users\n", reqUser.Role), "")
		return
	}

	// Decode form body.
	var newUser forms.UserSignup

	err := c.ShouldBindJSON(&newUser)
	if err != nil {
		utils.HandleAbort(c, http.StatusBadRequest, "could not unmarshal request body into 'forms.UserSignup'", err.Error())
		return
	}

	// Send sign up request to authentication service.
	conf := config.GetConfig()
	cookie, _ := c.Cookie("Authorization") // cookie is set (verified at the middleware)

	header := map[string]string{"Authorization": cookie}
	payload := struct {
		Id       string
		Password string
	}{Id: newUser.UId, Password: newUser.Password}

	response, err := utils.SendHTTPRequest(
		conf.GetString("endpoints.auth.signup.method"),
		conf.GetString("endpoints.auth.signup.method.base")+conf.GetString("endpoints.auth.signup.method.uri"),
		header,
		payload,
	)

	if (err != nil) || (response.StatusCode != http.StatusCreated) {
		utils.HandleAbort(c, http.StatusInternalServerError, "auth service not up or invalid request", err.Error())
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

// Update        godoc
// @Summary      Update a user
// @Description  When provided with a user object (complete or partial), the user specified in path is updated.
// @Tags         users
// @Produce      json
// @Param        user  body      forms.UserSignup  true  "User update form"
// @Success      200   {object}  models.User
// @Security 	 BearerToken
// @Router       /user [patch]
func (u UserController) Update(c *gin.Context) {

	user, exists := c.Get("user")
	if !exists {
		utils.HandleAbort(c, http.StatusNotFound, "user not found, even though is authenticated", "")
		return
	}

	userObj := user.(*models.User)

	// Decoding the request body.
	var userUpdate forms.UserUpdate

	err := c.ShouldBindJSON(&userUpdate)
	if err != nil {
		utils.HandleAbort(c, http.StatusBadRequest, "could not unmarshal request body into 'forms.UserUpdate'", err.Error())
		return
	}

	// Checking for password update and notify the auth service.
	if userUpdate.Password != nil {

		// Send sign up request to authentication service.
		conf := config.GetConfig()
		cookie, _ := c.Cookie("Authorization") // cookie is set (verified at the middleware)

		header := map[string]string{"Authorization": cookie}
		payload := struct {
			Password string
		}{Password: *userUpdate.Password}

		response, err := utils.SendHTTPRequest(
			conf.GetString("endpoints.auth.signup.method"),
			conf.GetString("endpoints.auth.signup.method.base")+conf.GetString("endpoints.auth.signup.method.uri"),
			header,
			payload,
		)

		if (err != nil) || (response.StatusCode != http.StatusCreated) {
			utils.HandleAbort(c, http.StatusInternalServerError, "auth service not up or invalid request", err.Error())
			return
		}

	}

	// Update the rest of the user.
	userObj, err = userModel.Update(userObj.UId, userUpdate)
	if err != nil {
		utils.HandleAbort(c, http.StatusInternalServerError, "could not update user information", err.Error())
		return
	}

	// Complete.
	c.JSON(http.StatusOK, gin.H{
		"user": userObj,
	})
	return

}
