package handler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Slemgrim/styx/service"
	"github.com/gorilla/mux"
)

type Upload struct {
	Service service.Attachment
}

func (u Upload) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	attachment, err := u.Service.Load(id)

	if err != nil {
		respondError(w, err)
		return
	}

	if attachment.IsUploaded {
		respondError(w, errors.New("File already uploaded"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondError(w, err)
		return
	}

	err = u.Service.ValidateFile(attachment, body)

	if err != nil {
		respondError(w, err)
		return
	}

	err = u.Service.Upload(attachment, body, nil)

	if err != nil {
		respondError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func respondError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, err)
}
