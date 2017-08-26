package handler

import (
	"net/http"

	"github.com/Slemgrim/styx/model"
	"github.com/Slemgrim/styx/service"
	"github.com/Slemgrim/jsonapi"
	//"github.com/gorilla/mux"

	validator "gopkg.in/go-playground/validator.v9"
	"log"
	"github.com/gorilla/mux"
	"fmt"
)

type Mail struct {
	Validator *validator.Validate
	Service  service.Mail

	JsonApi
}

func (a Mail) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.setMediaType(w)
	errors := a.validateJsonApiHeaders(r)
	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		jsonapi.MarshalErrors(w, errors)
		return
	}

	status := http.StatusOK
	var err error
	var payload model.Mail

	if r.Method == "POST" {
		status, payload, errors = a.createMail(r)
	} else {
		status, payload, errors = a.getMail(r)
	}

	w.WriteHeader(status)

	if(len(errors) > 0){
		jsonapi.MarshalErrors(w, errors)
		return
	}

	ptr := &payload
	err = jsonapi.MarshalPayload(w, ptr)

	if err != nil {
		log.Fatal(err)
		return
	}
}

func (a Mail) getMail(r *http.Request) (status int, payload model.Mail, errors []*jsonapi.ErrorObject) {
	vars := mux.Vars(r)
	id := vars["id"]

	payload, err := a.Service.Load(id)

	if err != nil {
		status = http.StatusNotFound
		errors =  append(errors, &jsonapi.ErrorObject{
			Title: "Not Found",
		})
		return
	}

	return
}

func (a Mail) createMail(r *http.Request) (status int, payload model.Mail, errors []*jsonapi.ErrorObject) {

	mail := new(model.Mail)
	if er := a.Unmarshal(r.Body, mail); er != nil {
		errors =  append(errors, er)
		status = http.StatusBadRequest
		return
	}

	err := a.Validator.Struct(*mail)
	if err != nil {
		errors = a.HandleValidationErrors(err)
		status = http.StatusBadRequest
		return
	}

	payload = *mail
	payload, err = a.Service.Create(*mail)

	if err != nil {
		fmt.Println(err)
		errors =  append(errors, a.Error(err))
		status = http.StatusBadRequest
		return
	}

	return
}
