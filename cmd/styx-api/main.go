package main

import (
	"fmt"
	"log"

	"github.com/fetzi/styx/config"
	"github.com/fetzi/styx/model"
	"github.com/fetzi/styx/queue"
	"github.com/fetzi/styx/resource"
	"github.com/fetzi/styx/storage"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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
		return
	}

	defer db.Close()

	queue, err := queue.NewConnection(config.Queue.Host, config.Queue.Port, config.Queue.Username, config.Queue.Password)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer queue.Close()

	router := gin.Default()
	api := api2go.NewAPIWithRouting(
		"",
		api2go.NewStaticResolver("/"),
		gingonic.New(router),
	)

	mailStatusStorage := storage.NewMailStatusStorage(db)

	api.AddResource(model.Mail{}, resource.MailResource{&mailStatusStorage, queue, config.Queue.QueueName})

	router.Run(fmt.Sprintf(":%d", config.HTTP.Port))

}
