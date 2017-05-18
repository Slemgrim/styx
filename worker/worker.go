package worker

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fetzi/styx/mailer"
	"github.com/fetzi/styx/model"
	"github.com/fetzi/styx/queue"
	"github.com/jinzhu/gorm"
)

// QueueWorker defines the queue specific information
type QueueWorker struct {
	Database        *gorm.DB
	QueueConnection *queue.Connection
	QueueName       string
	Mailer          *mailer.Mailer
}

type MailConsumer struct {
	channel chan model.Mail
	queue.MessageCallback
	Mailer *mailer.Mailer
}

// NewQueueWorker creates a new queue worker instance
func NewQueueWorker(database *gorm.DB, queueConnection *queue.Connection, queueName string, mailer *mailer.Mailer) *QueueWorker {
	return &QueueWorker{
		database,
		queueConnection,
		queueName,
		mailer,
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
		channel: queueToSMTP,
		Mailer:  worker.Mailer,
	})

	// wait for signal
	<-done
	fmt.Println("worker shutdown complete")
}

func (c MailConsumer) Execute(message queue.Message) {
	mail := model.Mail{}
	message.ParseFromJSON(&mail)

	err := c.Mailer.Send(mail)
	if err != nil {
		//Todo what to do when a queue entry can't be sent #19
	}

	message.Acknowledge()
}
