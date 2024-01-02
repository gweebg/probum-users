package seeder

import (
	"encoding/json"
	"github.com/gweebg/probum-users/database"
	"github.com/gweebg/probum-users/models"
	"github.com/gweebg/probum-users/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"os"
)

type Seeder struct {
	SeedPath string
	Db       *gorm.DB
}

func New(seedPath string) *Seeder {

	db := database.GetDatabase()

	return &Seeder{
		SeedPath: seedPath,
		Db:       db,
	}
}

type User struct {
	UId   string `json:"UId"`
	Email string `json:"Email"`

	Name string `json:"Name"`
	Role string `json:"Role"`

	Password string `json:"password"`
}

func (s *Seeder) Seed() {

	log.Println("seeding database")

	usersJson, err := os.ReadFile(s.SeedPath)
	utils.Check(err, "")

	var users []User

	err = json.Unmarshal(usersJson, &users)
	utils.Check(err, "")

	// Print the users' information
	for _, user := range users {

		log.Println(user.Password)

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		utils.Check(err, "cannot encrypt password '%v'\n", user.Password)

		dbUser := models.User{
			UId:      user.UId,
			Email:    user.Email,
			Name:     user.Name,
			Role:     user.Role,
			Password: hash,
		}

		s.Db.Create(&dbUser)
	}

	log.Println("finished seeding database")

}
