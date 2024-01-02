package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"gorm.io/gorm"

	"github.com/gweebg/probum-users/config"
	"github.com/gweebg/probum-users/database"
	"github.com/gweebg/probum-users/models"
	"github.com/gweebg/probum-users/seeder"
	"github.com/gweebg/probum-users/server"
)

func main() {

	environment := flag.String("e", "development", "")
	flag.Usage = func() {
		log.Println("usage: server -e {mode}")
		os.Exit(1)
	}

	flag.Parse()

	config.Init(*environment)

	database.ConnectDatabase()

	db := database.GetDatabase()

	if err := db.AutoMigrate(&models.User{}); err == nil && db.Migrator().HasTable(&models.User{}) {
		if err := db.First(&models.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			s := seeder.New("seeder/seed_data/users.json")
			s.Seed()
		}
	}

	server.Init()

}
