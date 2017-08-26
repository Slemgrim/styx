package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/Slemgrim/styx"
	"github.com/Slemgrim/styx/config"
	"github.com/Slemgrim/styx/handler"
	"github.com/Slemgrim/styx/model"
	"github.com/Slemgrim/styx/queue"
	"github.com/Slemgrim/styx/resource"
	"github.com/Slemgrim/styx/service"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

func main() {
	config, err := config.ReadConfig("config.json")

	if err != nil {
		log.Fatal(err)
	}

	/**
	 * Setup MongoDB
	 */

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
	defer session.Close()

	mResource := resource.MongoMail{Collection: db.C("mails")}
	aResource := resource.MongoAttachment{Collection: db.C("attachments")}

	/**
	 * Setup queueing system
	 */

	queue, err := queue.NewConnection(config.Queue.Host, config.Queue.Port, config.Queue.Username, config.Queue.Password)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer queue.Close()

	/**
	 * Add Custom Validators
	 */

	v := validator.New()
	v.RegisterStructValidation(model.ValidateBody, model.Body{})
	v.RegisterStructValidation(model.ValidateAddress, model.Address{})

	/**
	 * Register services
	 */

	aStore := styx.GetAttachmentStore(config.Files, db)
	aService := service.Attachment{Resource: aResource}

	mService := service.Mail{
		MailResource:       mResource,
		AttachmentResource: aResource,
		Connection:         queue,
	}

	aHandler := handler.Attachment{Validator: v, Service: aService}
	uHandler := handler.Upload{Service: aService, Store: aStore}
	mHandler := handler.Mail{Validator: v, Service: mService}

	/**
	 * Setup routing
	 */

	r := mux.NewRouter()
	r.Handle("/attachments", aHandler).Methods("POST")
	r.Handle("/attachments/{id}", aHandler).Methods("GET")
	r.Handle("/upload/{id}", uHandler).Methods("PUT")
	r.Handle("/mails", mHandler).Methods("POST")
	r.Handle("/mails/{id}", mHandler).Methods("GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":9999", nil))
}
