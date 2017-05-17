package main

import (
	"fmt"
	"log"

	"github.com/fetzi/styx/config"
	"github.com/fetzi/styx/queue"
	"github.com/fetzi/styx/worker"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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

	worker := worker.NewQueueWorker(db, queue, config.Queue.QueueName)

	worker.Start()
}
