package controllers

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gweebg/probum-users/config"
	"github.com/gweebg/probum-users/utils"
	"net/http"
	"time"

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
// @Security 	 BearerToken
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

	user, exists := c.Get("user")
	if !exists {
		utils.HandleAbort(c, http.StatusNotFound, "user not found, even though is authenticated", "")
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{"user": user.(*models.User)},
	)
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

	// Insert the new user into the database.
	user, err := userModel.Create(newUser)
	if err != nil {
		utils.HandleAbort(c, http.StatusInternalServerError, "could not add user to the database", err.Error())
		return
	}

	// Complete.
	c.JSON(
		http.StatusCreated,
		gin.H{"user": user},
	)
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

	// Update the rest of the user.
	userObj, err = userModel.Update(userObj.UId, userUpdate)
	if err != nil {
		utils.HandleAbort(c, http.StatusInternalServerError, "could not update user information", err.Error())
		return
	}

	// Complete.
	c.JSON(
		http.StatusOK,
		gin.H{"user": userObj},
	)
}

// Login         godoc
// @Summary      Authenticate user
// @Description  Authenticate user via email and password. Returns the user a jwt token if successful.
// @Tags         auth
// @Produce      json
// @Param        user  body      forms.UserLogin  true  "Login credentials"
// @Success      200   {object}  object
// @Router       /user/login [post]
func (u UserController) Login(c *gin.Context) {

	var login forms.UserLogin

	err := c.ShouldBindJSON(&login)
	if err != nil {
		utils.HandleAbort(c, http.StatusBadRequest, "could not unmarshal request body into 'forms.UserLogin'", err.Error())
		return
	}

	userId, err := userModel.CheckPassword(login)
	if err != nil {
		utils.HandleAbort(c, http.StatusUnauthorized, "invalid credentials", err.Error())
		return
	}

	tokenString, err := CreateToken(userId)
	if err != nil {
		utils.HandleAbort(c, http.StatusInternalServerError, "could not generate jwt", err.Error())
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{"token": tokenString},
	)

}

func CreateToken(userId string) (string, error) {

	c := config.GetConfig()

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": userId,
			"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // One Week
		},
	)

	jwtSecret := c.GetString("app.jwt-secret")
	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
