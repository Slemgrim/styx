package styx

import (
	"github.com/Slemgrim/gorage"
	"github.com/Slemgrim/gorage/meta"
	"github.com/Slemgrim/gorage/relation"
	"github.com/Slemgrim/gorage/storage"
	"github.com/Slemgrim/styx/config"
	"github.com/Slemgrim/styx/mailer"
	"github.com/Slemgrim/styx/model"
	"github.com/Slemgrim/styx/queue"
	"github.com/Slemgrim/styx/resource"
	"github.com/Slemgrim/styx/service"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2"
	"log"
)

type Styx struct {
	MailService       service.Mail
	AttachmentService service.Attachment
	Validator         *validator.Validate
	Session           *mgo.Session
	Queue             *queue.Connection
	Mailer            *mailer.Mailer
}

func New(config *config.Config) Styx {

	/**
	 * Setup MongoDB
	 */

	s := Styx{}

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    config.MongoDB.Address,
		Database: config.MongoDB.Database,
		Username: config.MongoDB.User,
		Password: config.MongoDB.Password,
	}

	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	db := session.DB("styx")
	if err != nil {
		panic(err)
	}

	s.Session = session

	mResource := resource.MongoMail{Collection: db.C("mails")}
	aResource := resource.MongoAttachment{Collection: db.C("attachments")}

	/**
	 * Setup queueing system
	 */

	queue, err := queue.NewConnection(config.Queue.Host, config.Queue.Port, config.Queue.Username, config.Queue.Password)

	if err != nil {
		log.Fatal(err)
	}
	s.Queue = queue

	/**
	 * Add Custom Validators
	 */

	s.Validator = validator.New()
	s.Validator.RegisterStructValidation(model.ValidateBody, model.Body{})
	s.Validator.RegisterStructValidation(model.ValidateAddress, model.Address{})

	/**
	 * Register services
	 */

	aStore := getAttachmentStore(config.Files, db)
	s.AttachmentService = service.Attachment{
		Resource: aResource,
		Store:    aStore,
	}

	s.MailService = service.Mail{
		MailResource:       mResource,
		AttachmentResource: aResource,
		Connection:         queue,
	}

	s.Mailer = mailer.New(config.SMTP)

	return s
}

func getAttachmentStore(config config.FilesConfig, database *mgo.Database) *gorage.Gorage {
	s := storage.Io{
		BasePath:   config.AttachmentPath,
		DirLength:  6,
		BufferSize: 1024,
	}

	r := relation.Mongo{Collection: database.C("attachment_relation")}
	m := meta.Mongo{Collection: database.C("attachment_meta")}

	gorage := gorage.NewGorage(s, r, m)

	return gorage
}

func (s *Styx) Close() {
	s.Session.Close()
	s.Queue.Close()
}
