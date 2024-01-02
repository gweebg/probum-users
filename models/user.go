package models

import (
	"github.com/gweebg/probum-users/database"
	"github.com/gweebg/probum-users/forms"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	UId   string `gorm:"primaryKey" json:"UId"`
	Email string `gorm:"unique;index" json:"Email"`

	Name string `gorm:"size:255" json:"Name"`
	Role string `gorm:"size:31" json:"Role"`

	Password []byte `gorm:"index" json:"password"`

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

	hash, err := bcrypt.GenerateFromPassword([]byte(userPayload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := User{
		UId:      userPayload.UId,
		Email:    userPayload.Email,
		Name:     userPayload.Name,
		Role:     userPayload.Role,
		Password: hash,
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u User) Update(userId string, userPayload forms.UserUpdate) (*User, error) {

	db := database.GetDatabase()
	tx := db.Begin()

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
	if userPayload.Password != nil {

		hash, err := bcrypt.GenerateFromPassword([]byte(*userPayload.Password), bcrypt.DefaultCost)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		user.Password = hash
	}

	tx.Save(user)
	tx.Commit()
	return user, nil
}

func (u User) CheckPassword(payload forms.UserLogin) (string, error) {

	db := database.GetDatabase()

	var user User
	if err := db.First(&user, "email", payload.Email).Error; err != nil {
		return "", err
	}

	err := bcrypt.CompareHashAndPassword(user.Password, []byte(payload.Password))
	if err != nil {
		return "", err
	}

	return user.UId, nil
}
