package resource

import (
	"errors"
	"net/http"

	"github.com/fetzi/styx/storage"
	"github.com/manyminds/api2go"
)

// MailStatusResource tba
type MailStatusResource struct {
	MailStatusStorage *storage.MailStatusStorage
}

func (s MailStatusResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	return &Response{}, api2go.NewHTTPError(errors.New("find one not allowed"), "find one not allowed", http.StatusBadRequest)
}

func (s MailStatusResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	mailStatus, err := s.MailStatusStorage.GetOne(ID)

	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, "not found", http.StatusNotFound)
	}

	return &Response{Res: mailStatus, Code: http.StatusOK}, nil
}

func (s MailStatusResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	return &Response{}, api2go.NewHTTPError(errors.New("delete not allowed"), "delete not allowed", http.StatusBadRequest)
}

func (s MailStatusResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	return &Response{}, api2go.NewHTTPError(errors.New("update not allowed"), "update not allowed", http.StatusBadRequest)
}
