package resource

import (
	"github.com/slemgrim/styx/model"
	"gopkg.in/mgo.v2"
	"log"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Mail interface {
	Create(mail model.Mail) (model.Mail, error)
	Read(id string) (model.Mail, error)
	Update(model.Mail) (model.Mail, error)
}

type MongoMail struct {
	Collection *mgo.Collection
}

func (a MongoMail) Create(mail model.Mail) (model.Mail, error) {
	err := a.Collection.Insert(&mail)

	if err != nil {
		log.Fatal(err)
	}

	return mail, nil
}

func (a MongoMail) Read(id string) (model.Mail, error) {
	mail := model.Mail{}
	err := a.Collection.Find(bson.M{"id": id, "deletedat": time.Time{}}).One(&mail)
	if err != nil {
		return mail, err
	}

	return mail, nil
}

func (a MongoMail) Update(mail model.Mail) (model.Mail, error) {
	err := a.Collection.Update(bson.M{"id": mail.ID, "deletedat": time.Time{}}, mail)
	if err != nil {
		return mail, err
	}

	return mail, nil
}
