package main

import (
	"fmt"
	"log"

	"github.com/Slemgrim/gorage"
	"github.com/Slemgrim/gorage/meta"
	"github.com/Slemgrim/gorage/relation"
	store "github.com/Slemgrim/gorage/storage"
	"github.com/fetzi/styx/config"
	"github.com/fetzi/styx/mailer"
	"github.com/fetzi/styx/queue"
	"github.com/fetzi/styx/storage"
	"github.com/fetzi/styx/worker"

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

	bodyStore := setupBodyStore(config, db)
	fmt.Println(bodyStore)

	queue, err := queue.NewConnection(config.Queue.Host, config.Queue.Port, config.Queue.Username, config.Queue.Password)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer queue.Close()

	mailer := mailer.NewMailer(config.SMTP, config.Files)
	mailStatusStorage := storage.NewMailStatusStorage(db)

	worker := worker.NewQueueWorker(db, queue, config.Queue.QueueName, mailer, bodyStore, mailStatusStorage)

	worker.Start()
}

func setupBodyStore(config *config.Config, db *gorm.DB) *gorage.Gorage {
	s := store.Io{
		BasePath:   config.Files.BodyPath,
		DirLength:  6,
		BufferSize: 1024,
	}

	r := relation.NewDb("mail_body_rel", db)
	m := meta.NewDb("mail_body", db)

	gorage := gorage.NewGorage(s, r, m)

	return gorage
}
