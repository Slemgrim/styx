package handler

import (
	"net/http"

	"github.com/fetzi/styx/model"
	"github.com/fetzi/styx/service"
	"github.com/google/jsonapi"
	"github.com/gorilla/mux"

	validator "gopkg.in/go-playground/validator.v9"
)

type Attachment struct {
	Validate *validator.Validate
	Service  service.Attachment

	JsonApi
}

func (a Attachment) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.setMediaType(w)
	if !a.validateJsonApiHeaders(w, r) {
		return
	}

	var err error

	if r.Method == "POST" {
		err = a.createAttachment(w, r)
	} else {
		err = a.getAttachment(w, r)
	}

	if err != nil {
		a.returnError(w, err, http.StatusInternalServerError)
		return
	}
}

func (a Attachment) getAttachment(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	attachment, err := a.Service.Load(id)

	if err != nil {
		a.returnError(w, err, http.StatusNotFound)
	}

	w.WriteHeader(http.StatusOK)
	err = jsonapi.MarshalOnePayload(w, &attachment)

	return err
}

func (a Attachment) createAttachment(w http.ResponseWriter, r *http.Request) error {
	attachment := new(model.Attachment)
	err := jsonapi.UnmarshalPayload(r.Body, attachment)
	if err != nil {
		a.returnError(w, err, http.StatusInternalServerError)
	}

	err = a.validateAttachment(w, *attachment)

	if err != nil {
		a.returnError(w, err, http.StatusInternalServerError)
	}

	newAttachment, err := a.Service.Create(*attachment)

	if err != nil {
		a.returnError(w, err, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	err = jsonapi.MarshalOnePayload(w, &newAttachment)

	return err
}

func (a Attachment) validateAttachment(w http.ResponseWriter, attachment model.Attachment) error {
	err := a.Validate.Struct(attachment)
	if err != nil {
		a.RetunValidationErros(w, err)
		return err
	}

	return nil
}
