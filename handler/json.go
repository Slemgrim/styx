package handler

import (
	"net/http"
	"strconv"
	"github.com/Slemgrim/jsonapi"
	"github.com/slemgrim/styx"
	"io"
	"fmt"
	"gopkg.in/go-playground/validator.v9"
)

type JsonApi struct{}

func (h *JsonApi) validateJsonApiHeaders(r *http.Request) []*jsonapi.ErrorObject {
	var errors []*jsonapi.ErrorObject

	if r.Header.Get("Content-Type") != jsonapi.MediaType {
		errors = append(errors, &jsonapi.ErrorObject{
			Title:  "Unsupported Content Type",
			Detail: "Given Content Type was not:" + jsonapi.MediaType,
			Code: styx.ErrorWrongContentType,
		})
	}

	if r.Header.Get("Accept") != jsonapi.MediaType {
		errors = append(errors, &jsonapi.ErrorObject{
			Title:  "Response not supported",
			Detail: "Client must support " + jsonapi.MediaType,
			Code: styx.ErrorNotAccepted,
		})
	}

	return errors
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

func (h *JsonApi) Unmarshal(r io.Reader, model interface{}) (error *jsonapi.ErrorObject) {
	err := jsonapi.UnmarshalPayload(r, model)
	if err != nil {
		fmt.Println(err)
		error = new(jsonapi.ErrorObject)
		error.Code = styx.ErrorInvalidJson
		error.Title = "Can't unmarshal json"
		error.Detail = err.Error()
		return
	}

	return nil
}


func (a Attachment) HandleValidationErrors(errors error) (jsonErrors []*jsonapi.ErrorObject){

	for _, err := range errors.(validator.ValidationErrors) {
		jsonErrors = append(jsonErrors, &jsonapi.ErrorObject{
			Title: "Validation error for field: " + err.Field(),
			Meta: &map[string]interface{}{
				"field": err.StructField(),
				"tag":   err.Tag()},
			Code: styx.ErrorValidation,
		})
	}

	return
}