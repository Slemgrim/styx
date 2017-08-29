package worker

import (
	"fmt"
	"github.com/Slemgrim/styx"
	"github.com/Slemgrim/styx/mailer"
	"github.com/Slemgrim/styx/model"
	"github.com/Slemgrim/styx/queue"
	"github.com/Slemgrim/styx/service"
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Worker struct {
	session           *mgo.Session
	queue             *queue.Connection
	mailer            *mailer.Mailer
	mailService       service.Mail
	attachmentService service.Attachment
}

func New(styx *styx.Styx) *Worker {

	worker := &Worker{
		styx.Session,
		styx.Queue,
		styx.Mailer,
		styx.MailService,
		styx.AttachmentService,
	}
	return worker
}

func (w *Worker) Start() error {
	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	channel, err := w.queue.Channel()

	if err != nil {
		return err
	}

	defer channel.Close()

	go func() {
		<-signals
		fmt.Println("recieved shutdown signal")
		done <- true
	}()

	q, err := channel.DeclareQueue("mails", false, false, false, false)

	if err != nil {
		return err
	}

	channel.Prefetch(20)
	channel.Consume(q, "styx-consumer", mailConsumer{
		mailer:            w.mailer,
		mailService:       w.mailService,
		attachmentService: w.attachmentService,
	})

	<-done
	fmt.Println("worker shutdown complete")

	return nil
}

func (w *Worker) Stop() {

}

type mailConsumer struct {
	mailer            *mailer.Mailer
	mailService       service.Mail
	attachmentService service.Attachment
}

func (c mailConsumer) Execute(message queue.Message) {
	mail := model.Mail{}
	message.ParseFromJSON(&mail)

	err := c.mailer.Send(mail)
	if err != nil {
		log.Fatal(err)
	}

	c.mailService.MarkAsSent(mail.ID)

	message.Acknowledge()
}
