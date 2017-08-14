package resource

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/slemgrim/styx/model"
	"github.com/slemgrim/styx/queue"
	"github.com/slemgrim/styx/storage"
	"github.com/google/uuid"
	"github.com/manyminds/api2go"
)

type MailResource struct {
	MailStatusStorage *storage.MailStatusStorage
	QueueConnection   *queue.Connection
	QueueName         string
}

func (s MailResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	mail, ok := obj.(model.Mail)

	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	mail.ID = uuid.New().String()

	var from []string
	var to []string

	for _, client := range mail.Clients {
		switch client.Type {
		case "to":
			to = append(to, client.Email)
		case "from":
			from = append(from, client.Email)
		}
	}

	mailStatus := model.MailStatus{
		MailID:  mail.ID,
		Subject: mail.Subject,
		From:    strings.Join(from, ", "),
		To:      strings.Join(to, ", "),
		Created: time.Now().Unix(),
		Sent:    0,
	}

	s.MailStatusStorage.Insert(mailStatus)

	err := s.publishToQueue(mail)

	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, "Unable to publish mail entry queue", http.StatusInternalServerError)
	}

	return &Response{Res: mailStatus, Code: http.StatusCreated}, nil
}

func (s MailResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	mailStatus, err := s.MailStatusStorage.GetOne(ID)

	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, "not found", http.StatusNotFound)
	}

	return &Response{Res: mailStatus, Code: http.StatusOK}, nil
}

func (s MailResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	return &Response{}, api2go.NewHTTPError(errors.New("delete not allowed"), "delete not allowed", http.StatusBadRequest)
}

func (s MailResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	return &Response{}, api2go.NewHTTPError(errors.New("update not allowed"), "update not allowed", http.StatusBadRequest)
}

func (s MailResource) publishToQueue(mail model.Mail) error {
	channel, err := s.QueueConnection.Channel()

	if err != nil {
		log.Fatal(err)
		return errors.New("Unable to open queue channel")
	}

	defer channel.Close()

	queue, err := channel.DeclareQueue(s.QueueName, false, false, false, false)

	if err != nil {
		log.Fatal(err)
		return errors.New("Unable to declare queue")
	}

	err = channel.PublishAsJSON(queue, mail)

	if err != nil {
		return errors.New("Unable to publish e-mail to queue")
	}

	return nil
}
