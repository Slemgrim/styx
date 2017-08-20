package service

import (
	"time"

	"github.com/slemgrim/styx/model"
	"github.com/slemgrim/styx/resource"
	"github.com/google/uuid"
	"fmt"
	"errors"
)

/*
	Service for handling mails
*/
type Mail struct {
	MailResource resource.Mail
	AttachmentResource resource.Attachment
}

/*
	Crete a new mails
*/
func (s Mail) Create(mail model.Mail) (model.Mail, error) {
	var err error

	var attachments []*model.Attachment
	for _, attachment := range mail.Attachments {
		a, err := s.AttachmentResource.Read(attachment.ID)
		if err != nil{
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

	return mail, err
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