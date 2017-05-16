package styx

import (
	"fmt"
	"log"

	"github.com/fetzi/styx/model"
	"github.com/fetzi/styx/queue"
)

type MailConsumer struct {
	channel chan model.Mail
	queue.MessageCallback
}

func StartWorker(queueConnection *queue.Connection, queueName string) {
	go func() {
		//signals := make(chan os.Signal, 1)
		//done := make(chan bool, 1)

		queueToSMTP := make(chan model.Mail, 20)

		//signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

		channel, err := queueConnection.Channel()

		if err != nil {
			log.Fatal(err)
			return
		}

		/*go func() {
			<-signals
			channel.Close()
			done <- true
		}()*/

		q, err := channel.DeclareQueue(queueName, false, false, false, false)

		if err != nil {
			log.Fatal(err)
			return
		}

		channel.Consume(q, "styx-consumer", MailConsumer{channel: queueToSMTP})

		// wait for signal
		//<-done
	}()
}

func (c MailConsumer) Execute(message queue.Message) {
	mail := model.Mail{}
	message.ParseFromJSON(&mail)

	c.channel <- mail
	fmt.Print("%+v\n", mail)

	message.Acknowledge()
}
