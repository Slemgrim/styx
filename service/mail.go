package service

import (
	"time"

	"github.com/slemgrim/styx/model"
	"github.com/slemgrim/styx/resource"
	"github.com/google/uuid"
)

/*
	Service for handling mails
*/
type Mail struct {
	Resource resource.Mail
}

/*
	Crete a new mails
*/
func (s Mail) Create(mail model.Mail) (model.Mail, error) {
	var err error
	mail.ID = uuid.New().String()
	mail.CreatedAt = time.Now()
	mail.SentAt = time.Time{}
	mail.DeletedAt = time.Time{}

	mail, err = s.Resource.Create(mail)

	return mail, err
}

/*
	Load a mail by its id
*/
func (s Mail) Load(id string) (model.Mail, error) {

	mail, err := s.Resource.Read(id)

	if err != nil {
		return model.Mail{}, err
	}

	return mail, nil
}