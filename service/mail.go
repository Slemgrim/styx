package service

import (
	"time"

	"errors"
	"fmt"
	"github.com/Slemgrim/styx/model"
	"github.com/Slemgrim/styx/queue"
	"github.com/Slemgrim/styx/resource"
	"github.com/google/uuid"
	"log"
)

/*
	Service for handling mails
*/
type Mail struct {
	MailResource       resource.Mail
	AttachmentResource resource.Attachment
	Connection         *queue.Connection
}

/*
	Crete a new mails
*/
func (s Mail) Create(mail model.Mail) (model.Mail, error) {
	var err error

	var attachments []*model.Attachment
	for _, attachment := range mail.Attachments {
		a, err := s.AttachmentResource.Read(attachment.ID)
		if err != nil {
			return mail, err
		}

		if !a.IsUploaded {
			return mail, errors.New(fmt.Sprintf("Attachment %s was not uploaded yet", attachment.ID))
		}

		attachments = append(attachments, &a)
	}

	mail.Attachments = attachments
	mail.ID = uuid.New().String()
	mail.CreatedAt = time.Now()
	mail.SentAt = time.Time{}
	mail.DeletedAt = time.Time{}

	mail, err = s.MailResource.Create(mail)

	if err != nil {
		return mail, err
	}

	s.enqueue(mail)

	return mail, nil
}

/*
	Load a mail by its id
*/
func (s Mail) Load(id string) (model.Mail, error) {

	mail, err := s.MailResource.Read(id)

	if err != nil {
		return model.Mail{}, err
	}

	return mail, nil
}

func (s Mail) enqueue(mail model.Mail) {
	channel, err := s.Connection.Channel()

	if err != nil {
		log.Fatal(err)
	}

	defer channel.Close()

	queue, err := channel.DeclareQueue("mails", false, false, false, false)

	if err != nil {
		log.Fatal(err)
	}

	channel.PublishAsJSON(queue, mail)
}

func (s Mail) MarkAsSent(id string) error {

	mail, err := s.MailResource.Read(id)

	if err != nil {
		return err
	}

	mail.SentAt = time.Now()

	_, err = s.MailResource.Update(mail)

	if err != nil {
		return err
	}

	return nil
}
