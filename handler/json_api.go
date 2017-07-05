package handler

import (
	"net/http"
	"strconv"

	validator "gopkg.in/go-playground/validator.v9"

	"github.com/google/jsonapi"
)

type JsonApi struct{}

func (h *JsonApi) validateJsonApiHeaders(w http.ResponseWriter, r *http.Request) bool {
	var errors []*jsonapi.ErrorObject

	if r.Header.Get("Content-Type") != jsonapi.MediaType {
		errors = append(errors, &jsonapi.ErrorObject{
			Title:  "Unsupported Content Type",
			Detail: "Given Content Type was not:" + jsonapi.MediaType,
			Status: "400",
		})
	}

	if r.Header.Get("Accept") != jsonapi.MediaType {
		errors = append(errors, &jsonapi.ErrorObject{
			Title:  "Response not supported",
			Detail: "Client must support " + jsonapi.MediaType + "type",
			Status: "400",
		})
	}

	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)

		jsonapi.MarshalErrors(w, errors)
		return false
	}

	return true
}

func (h *JsonApi) setMediaType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", jsonapi.MediaType)
}

func (h *JsonApi) returnError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
		Title:  err.Error(),
		Status: strconv.Itoa(status),
	}})
}

func (h *JsonApi) returnErrors(w http.ResponseWriter, errors []error, status int) {
	w.WriteHeader(status)
	var e []*jsonapi.ErrorObject

	for _, err := range errors {
		e = append(e, &jsonapi.ErrorObject{
			Title:  err.Error(),
			Status: strconv.Itoa(status),
		})
	}
	jsonapi.MarshalErrors(w, e)
}

func (h *JsonApi) RetunValidationErros(w http.ResponseWriter, errors error) {
	w.WriteHeader(http.StatusBadRequest)
	var e []*jsonapi.ErrorObject

	for _, err := range errors.(validator.ValidationErrors) {
		e = append(e, &jsonapi.ErrorObject{
			Title:  "Validation error",
			Detail: "Validation error for field: " + err.Field(),
			Meta: &map[string]interface{}{
				"field": err.StructField(),
				"tag":   err.Tag()},
			Status: strconv.Itoa(http.StatusBadRequest),
		})
	}

	jsonapi.MarshalErrors(w, e)
}
