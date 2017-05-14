package main

import (
	"log"

	"github.com/gin-gonic/gin"
    "github.com/manyminds/api2go"
    "github.com/manyminds/api2go-adapter/gingonic"
	"github.com/fetzi/styx/config"
)

func main() {
	config, err := config.ReadConfig("config.json")

	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(config.Storage.Driver, config.Storage.Config)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	router := gin.Default()
	api := api2go.NewAPIWithRouting(
		"api",
		api2Go.NewStaticResolver("/"),
		gingonic.New(router)
	)

	mailStatusStorage = storage.NewMailStatusStorage()
}
