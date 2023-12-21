package seeder

import (
	"encoding/json"
	"github.com/gweebg/probum-users/database"
	"github.com/gweebg/probum-users/models"
	"github.com/gweebg/probum-users/utils"
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

func (s *Seeder) Seed() {

	log.Println("seeding database")

	usersJson, err := os.ReadFile(s.SeedPath)
	utils.Check(err, "")

	var users []models.User

	err = json.Unmarshal(usersJson, &users)
	utils.Check(err, "")

	// Print the users' information
	for _, user := range users {
		s.Db.Create(&user)
	}

	log.Println("finished seeding database")

}
