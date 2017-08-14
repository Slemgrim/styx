package main

import (
	"fmt"
	"log"

	"github.com/slemgrim/styx/config"
	"github.com/slemgrim/styx/mailer"
	"github.com/slemgrim/styx/queue"
	"github.com/slemgrim/styx/worker"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	fmt.Println("starting styx worker")
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

	mailer := mailer.NewMailer(config.SMTP, config.Attachments)
	worker := worker.NewQueueWorker(db, queue, config.Queue.QueueName, mailer)

	worker.Start()
}
