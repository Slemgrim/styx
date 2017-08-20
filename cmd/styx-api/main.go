package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/slemgrim/styx"
	"github.com/slemgrim/styx/config"
	"github.com/slemgrim/styx/handler"
	"github.com/slemgrim/styx/resource"
	"github.com/slemgrim/styx/service"
	"github.com/gorilla/mux"
	"github.com/slemgrim/styx/model"
)

func main() {
	config, err := config.ReadConfig("config.json")

	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(config.Storage.Driver, config.Storage.Config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	v := validator.New()
	v.RegisterStructValidation(model.ValidateBody, model.Body{})
	v.RegisterStructValidation(model.ValidateAddress, model.Address{})

	aStore := styx.GetAttachmentStore(config.Files, db)

	aResource := resource.DbAttachment{DB: db}
	aResource.Init()
	aService := service.Attachment{Resource: aResource}

	mResource := resource.DbMail{DB: db}
	mResource.Init()
	mService := service.Mail{Resource: mResource}

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
