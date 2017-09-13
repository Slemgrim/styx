package main

import (
	"log"
	"net/http"

	"github.com/Slemgrim/styx"
	"github.com/Slemgrim/styx/config"
	"github.com/Slemgrim/styx/handler"
	"github.com/gorilla/mux"
)

func main() {
	config, err := config.ReadConfig("config.json")

	if err != nil {
		log.Fatal(err)
	}

	styx := styx.New(config)
	defer styx.Close()

	aHandler := handler.Attachment{Validator: styx.Validator, Service: styx.AttachmentService}
	uHandler := handler.Upload{Service: styx.AttachmentService}
	mHandler := handler.Mail{Validator: styx.Validator, Service: styx.MailService}

	r := mux.NewRouter()
	r.Handle("/attachments", aHandler).Methods("POST")
	r.Handle("/attachments/{id}", aHandler).Methods("GET")
	r.Handle("/upload/{id}", uHandler).Methods("PUT")
	r.Handle("/mails", mHandler).Methods("POST")
	r.Handle("/mails/{id}", mHandler).Methods("GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":"+config.HTTP.Port, nil))
}
