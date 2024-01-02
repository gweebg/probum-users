package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gweebg/probum-users/config"
	"github.com/gweebg/probum-users/models"
	"net/http"
	"strings"
	"time"
)

func RequireAuth(c *gin.Context) {

	conf := config.GetConfig()

	userModel := new(models.User)

	// Get the cookie from the request.
	bearer := c.GetHeader("Authorization")
	if bearer == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "authorization token not set",
		})
		return
	}

	bearer = strings.Split(bearer, " ")[1]

	// Decode and validate the cookie.
	token, err := jwt.Parse(bearer, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected jwt signing method: " + token.Method.Alg())
		}
		return []byte(conf.GetString("app.jwt-secret")), nil

	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "token is not valid",
			"error":   err.Error(),
		})
		return
	}
	claims := token.Claims.(jwt.MapClaims)

	// Check the expiration.
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "token is expired",
		})
		return
	}

	// Find the user with the token subject.
	user, err := userModel.Get(claims["sub"].(string)) // claims["sub"] contains the user id (UId field)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "no matching user to token",
		})
		return
	}

	// Attach to request & Continue
	c.Set("user", user)
	c.Next()
}
