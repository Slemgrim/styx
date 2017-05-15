package main

import (
	"log"

	"github.com/fetzi/styx/config"
	"github.com/fetzi/styx/model"
	"github.com/fetzi/styx/resource"
	"github.com/fetzi/styx/storage"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/manyminds/api2go"
	"github.com/manyminds/api2go-adapter/gingonic"
	gin "gopkg.in/gin-gonic/gin.v1"
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
		api2go.NewStaticResolver("/"),
		gingonic.New(router),
	)

	mailStatusStorage := storage.NewMailStatusStorage(db)

	api.AddResource(model.Mail{}, resource.MailResource{&mailStatusStorage})

	router.Run(":9999")
}
