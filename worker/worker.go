package worker

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Slemgrim/gorage"
	"github.com/fetzi/styx/mailer"
	"github.com/fetzi/styx/model"
	"github.com/fetzi/styx/queue"
	"github.com/fetzi/styx/storage"

	"github.com/jinzhu/gorm"
)

// QueueWorker defines the queue specific information
type QueueWorker struct {
	Database        *gorm.DB
	QueueConnection *queue.Connection
	QueueName       string
	Mailer          *mailer.Mailer
	BodyStore       *gorage.Gorage
	StatusStorage   storage.MailStatusStorage
}

type MailConsumer struct {
	channel chan model.Mail
	queue.MessageCallback
	mailer        *mailer.Mailer
	bodyStore     *gorage.Gorage
	statusStorage storage.MailStatusStorage
	database      *gorm.DB
}

// NewQueueWorker creates a new queue worker instance
func NewQueueWorker(
	database *gorm.DB,
	queueConnection *queue.Connection,
	queueName string,
	mailer *mailer.Mailer,
	bodyStore *gorage.Gorage,
	statusStorage storage.MailStatusStorage) *QueueWorker {
	return &QueueWorker{
		database,
		queueConnection,
		queueName,
		mailer,
		bodyStore,
		statusStorage,
	}
}

// Start starts the worker execution
func (worker *QueueWorker) Start() {
	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	queueToSMTP := make(chan model.Mail, 20)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	channel, err := worker.QueueConnection.Channel()

	if err != nil {
		log.Fatal(err)
		return
	}

	defer channel.Close()

	go func() {
		<-signals
		fmt.Println("recieved shutdown signal")
		done <- true
	}()

	q, err := channel.DeclareQueue(worker.QueueName, false, false, false, false)

	if err != nil {
		log.Fatal(err)
		return
	}

	channel.Prefetch(20)
	channel.Consume(q, "styx-consumer", MailConsumer{
		channel:       queueToSMTP,
		mailer:        worker.Mailer,
		bodyStore:     worker.BodyStore,
		statusStorage: worker.StatusStorage,
		database:      worker.Database,
	})

	// wait for signal
	<-done
	fmt.Println("worker shutdown complete")
}

func (c MailConsumer) Execute(message queue.Message) {
	mail := model.Mail{}
	message.ParseFromJSON(&mail)

	err := c.mailer.Send(mail)
	if err != nil {
		fmt.Println(err)
		//Todo what to do when a queue entry can't be sent #19
	}

	context := make(map[string]string)
	context["ID"] = mail.ID
	context["context"] = mail.Context
	context["type"] = "HTML"

	bodyHTMLFile, err := c.bodyStore.Save(mail.ID, []byte(mail.Body.HTML), context)

	if err != nil {
		fmt.Println(err)
		return
	}

	context["type"] = "Plain"

	bodyPlainFile, err := c.bodyStore.Save(mail.ID, []byte(mail.Body.Plain), context)

	if err != nil {
		fmt.Println(err)
		return
	}

	status, err := c.statusStorage.GetOne(mail.ID)

	if err != nil {
		fmt.Println(err)
		return
	}

	status.BodyHtml = bodyHTMLFile.ID
	status.BodyPlain = bodyPlainFile.ID
	status.Sent = 1

	c.database.Save(&status)

	message.Acknowledge()
}
