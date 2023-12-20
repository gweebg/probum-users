package main

import (
	"flag"
	"github.com/gweebg/probum-users/config"
	"github.com/gweebg/probum-users/database"
	"log"
	"os"
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
	database.MigrateDatabase()

	server.Init()

}
