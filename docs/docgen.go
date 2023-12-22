package docs

// @title           Probum User Manager
// @version         1.0
// @description     A user management service for the application Probum.
// @contact.name   Guilherme
// @contact.url    https://github.com/gweebg
// @license.name  MIT
// @license.url   https://mit-license.org/
// @host      localhost:3000
// @BasePath  /api/v1

// GetUser       godoc
// @Summary      Retrieve a user from the database,
// @Description  Retrieves the user with the specified school id, as a json object.
// @Tags         users
// @Produce      json
// @Param        id  path      string  true  "search book by isbn"
// @Success      200   {object}  models.User
// @Router       /user/{id} [get]
