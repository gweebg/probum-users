package database

import (
	"fmt"

	"github.com/gweebg/probum-users/config"
	"github.com/gweebg/probum-users/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDatabase() {
	var err error

	db, err = gorm.Open(postgres.Open(dsnFromConfig()), &gorm.Config{})
	utils.Check(err, "")
}

func dsnFromConfig() string {

	// DSN="host=localhost user=guilherme password=users dbname=users port=5432 sslmode=disable TimeZone=Europe/Lisbon"
	c := config.GetConfig()

	host := c.GetString("db.host")
	user := c.GetString("db.user")
	password := c.GetString("db.password")
	dbname := c.GetString("db.dbname")
	port := c.GetInt("db.port")
	tz := c.GetString("db.tz")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		host, user, password, dbname, port, tz,
	)

	return dsn
}

func GetDatabase() *gorm.DB {
	return db
}
