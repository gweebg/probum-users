package models

import (
	"gorm.io/gorm"

	"github.com/gweebg/probum-users/config"
	"github.com/gweebg/probum-users/database"
	"github.com/gweebg/probum-users/forms"
	"github.com/gweebg/probum-users/utils"
)

type User struct {
	gorm.Model

	Id    string `gorm:"primaryKey"`
	Email string `gorm:"unique;index"`

	Name string `gorm:"size:255"`
	Role string `gorm:"size:31"`
}

func (u User) GetUserById(id string) (*User, error) {

	db := database.GetDatabase()

	var user User

	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u User) Signup(userPayload forms.UserSignup) (*User, error) {

	c := config.GetConfig()
	db := database.GetDatabase()

	user := User{
		Id:    userPayload.Id,
		Email: userPayload.Email,
		Name:  userPayload.Name,
		Role:  userPayload.Role,
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	headers := map[string]string{}
	authUser := forms.AuthUser{
		Id:       userPayload.Id,
		Password: userPayload.Password,
	}

	_, err := utils.SendHTTPRequest(
		c.GetString("endpoints.auth.signup.method"),
		c.GetString("endpoints.auth.base")+c.GetString("endpoints.auth.signup.uri"),
		headers, authUser,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u User) Update(userId string, userPayload forms.UserUpdate) (*User, error) {

	//c := config.GetConfig()
	db := database.GetDatabase()

	user, err := u.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	if userPayload.Name != nil {
		user.Name = *userPayload.Name
	}
	if userPayload.Email != nil {
		user.Email = *userPayload.Email
	}

	if userPayload.Password != nil {
		// todo: if password is != nil then send request to auth service
	}

	db.Save(user)
	return user, nil
}
