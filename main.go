package main

import (
	"flag"
	"github.com/gweebg/probum-users/models"
	"github.com/gweebg/probum-users/utils"
	"log"
	"os"

	"github.com/gweebg/probum-users/config"
	"github.com/gweebg/probum-users/database"
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
	err := db.AutoMigrate(&models.User{})
	utils.Check(err, "")

	server.Init()

}
