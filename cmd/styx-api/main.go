package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/fetzi/styx"
	"github.com/fetzi/styx/config"
	"github.com/fetzi/styx/handler"
	"github.com/fetzi/styx/resource"
	"github.com/fetzi/styx/service"
	"github.com/gorilla/mux"
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

	aStore := styx.GetAttachmentStore(config.Files, db)

	aResource := resource.DbAttachment{DB: db}
	aResource.Init()
	aService := service.Attachment{Resource: aResource}

	aHandler := handler.Attachment{Validate: v, Service: aService}
	uHandler := handler.Upload{Service: aService, Store: aStore}

	r := mux.NewRouter()
	r.Handle("/attachments", aHandler).Methods("POST")
	r.Handle("/attachments/{id}", aHandler).Methods("GET")

	r.Handle("/upload/{id}", uHandler).Methods("POST")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":9999", nil))
}
