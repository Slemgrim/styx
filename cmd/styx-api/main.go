package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/slemgrim/styx"
	"github.com/slemgrim/styx/config"
	"github.com/slemgrim/styx/handler"
	"github.com/slemgrim/styx/resource"
	"github.com/slemgrim/styx/service"
	"github.com/gorilla/mux"
	"github.com/slemgrim/styx/model"
	"github.com/jinzhu/gorm"
	"gopkg.in/mgo.v2"
)

func main() {
	config, err := config.ReadConfig("config.json")

	if err != nil {
		log.Fatal(err)
	}

	//Gorm needed for gorage package. We should change this in the future
	gorm, err := gorm.Open(config.Storage.Driver, config.Storage.Config)
	if err != nil {
		log.Fatal(err)
	}
	defer gorm.Close()

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    config.MongoDB.Address,
		Database: config.MongoDB.Database,
		Username: config.MongoDB.User,
		Password: config.MongoDB.Password,
	}

	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	v := validator.New()
	v.RegisterStructValidation(model.ValidateBody, model.Body{})
	v.RegisterStructValidation(model.ValidateAddress, model.Address{})

	mResource := resource.MongoMail{Collection: session.DB("styx").C("mails")}
	aResource := resource.MongoAttachment{Collection: session.DB("styx").C("attachments")}

	aStore := styx.GetAttachmentStore(config.Files, gorm)
	aService := service.Attachment{Resource: aResource}
	mService := service.Mail{
		MailResource: mResource,
		AttachmentResource: aResource,
	}

	aHandler := handler.Attachment{Validator: v, Service: aService}
	uHandler := handler.Upload{Service: aService, Store: aStore}
	mHandler := handler.Mail{Validator: v, Service: mService}

	r := mux.NewRouter()
	r.Handle("/attachments", aHandler).Methods("POST")
	r.Handle("/attachments/{id}", aHandler).Methods("GET")
	r.Handle("/upload/{id}", uHandler).Methods("PUT")
	r.Handle("/mails", mHandler).Methods("POST")
	r.Handle("/mails/{id}", mHandler).Methods("GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":9999", nil))
}
