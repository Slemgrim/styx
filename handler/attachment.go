package handler

import (
	"net/http"

	"github.com/Slemgrim/styx/model"
	"github.com/Slemgrim/styx/service"
	"github.com/Slemgrim/jsonapi"
	"github.com/gorilla/mux"

	validator "gopkg.in/go-playground/validator.v9"
	"log"
)

type Attachment struct {
	Validator *validator.Validate
	Service  service.Attachment

	JsonApi
}

func (a Attachment) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.setMediaType(w)
	errors := a.validateJsonApiHeaders(r)
	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		jsonapi.MarshalErrors(w, errors)
		return
	}

	status := http.StatusOK
	var err error
	var payload model.Attachment

	if r.Method == "POST" {
		status, payload, errors = a.createAttachment(r)
	} else {
		status, payload, errors = a.getAttachment(r)
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

func (a Attachment) getAttachment(r *http.Request) (status int, payload model.Attachment, errors []*jsonapi.ErrorObject) {
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

func (a Attachment) createAttachment(r *http.Request) (status int, payload model.Attachment, errors []*jsonapi.ErrorObject) {

	attachment := new(model.Attachment)
	if er := a.Unmarshal(r.Body, attachment); er != nil {
		errors =  append(errors, er)
		status = http.StatusBadRequest
		return
	}

	err := a.Validator.Struct(*attachment)
	if err != nil {
		errors = a.HandleValidationErrors(err)
		status = http.StatusBadRequest
		return
	}

	payload, err = a.Service.Create(*attachment)

	if err != nil {
		return
	}

	return
}
