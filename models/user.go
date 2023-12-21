package models

import (
	"github.com/gweebg/probum-users/database"
	"github.com/gweebg/probum-users/forms"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	UId   string `gorm:"primaryKey" json:"UId"`
	Email string `gorm:"unique;index" json:"Email"`

	Name string `gorm:"size:255" json:"Name"`
	Role string `gorm:"size:31" json:"Role"`

	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdatedAt int64 `gorm:"autoUpdateTime:milli"`
}

func (u User) Get(id string) (*User, error) {

	db := database.GetDatabase()

	var user User

	if err := db.First(&user, "uid", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u User) Create(userPayload forms.UserSignup) (*User, error) {

	db := database.GetDatabase()

	user := User{
		UId:   userPayload.UId,
		Email: userPayload.Email,
		Name:  userPayload.Name,
		Role:  userPayload.Role,
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	// todo: move microservice call to controller instead of model
	// headers := map[string]string{}
	// authUser := forms.AuthUser{
	// 	UId:      userPayload.UId,
	// 	Password: userPayload.Password,
	// }
	//
	// _, err := utils.SendHTTPRequest(
	// 	c.GetString("endpoints.auth.signup.method"),
	// 	c.GetString("endpoints.auth.base")+c.GetString("endpoints.auth.signup.uri"),
	// 	headers, authUser,
	// )
	// if err != nil {
	// 	return nil, err
	// }

	return &user, nil
}

func (u User) Update(userId string, userPayload forms.UserUpdate) (*User, error) {

	db := database.GetDatabase()

	user, err := u.Get(userId)
	if err != nil {
		return nil, err
	}

	if userPayload.Name != nil {
		user.Name = *userPayload.Name
	}
	if userPayload.Email != nil {
		user.Email = *userPayload.Email
	}

	db.Save(user)
	return user, nil
}
