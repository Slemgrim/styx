package resource

import (
	"github.com/slemgrim/styx/model"
	"gopkg.in/mgo.v2"
	"log"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Attachment interface {
	Create(attachment model.Attachment) (model.Attachment, error)
	Read(id string) (model.Attachment, error)
	Update(model.Attachment) (model.Attachment, error)
}

type MongoAttachment struct {
	Collection *mgo.Collection
}

func (a MongoAttachment) Create(attachment model.Attachment) (model.Attachment, error) {
	err := a.Collection.Insert(&attachment)

	if err != nil {
		log.Fatal(err)
	}

	return attachment, nil
}

func (a MongoAttachment) Read(id string) (model.Attachment, error) {

	attachment := model.Attachment{}
	err := a.Collection.Find(bson.M{"id": id, "deletedat": time.Time{}}).One(&attachment)
	if err != nil {
		return attachment, err
	}

	return attachment, nil
}

func (a MongoAttachment) Update(attachment model.Attachment) (model.Attachment, error) {
	err := a.Collection.Update(bson.M{"id": attachment.ID, "deletedat": time.Time{}}, attachment)
	if err != nil {
		return attachment, err
	}

	return attachment, nil
}