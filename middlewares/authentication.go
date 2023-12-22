package middlewares

import (
	"github.com/gweebg/probum-users/config"
	"github.com/gweebg/probum-users/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"errors"
	"net/http"
	"time"
)

func RequireAuth(c *gin.Context) {

	conf := config.GetConfig()

	userModel := new(models.User)

	// Get the cookie from the request.
	cookie, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "authorization cookie not set",
		})
		return
	}

	// Decode and validate the cookie.
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected jwt signing method: " + token.Method.Alg())
		}
		return []byte(conf.GetString("jwt-secret")), nil

	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "token is not valid",
		})
		return
	}

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
