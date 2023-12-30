package server

import (
	"github.com/gweebg/probum-users/config"
	"github.com/gweebg/probum-users/utils"
)

func Init() {

	conf := config.GetConfig()
	router := NewRouter()

	err := router.Run(conf.GetString("app.listen"))
	utils.Check(err, "")
}
