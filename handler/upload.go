package handler

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Slemgrim/gorage"
	"github.com/Slemgrim/styx/model"
	"github.com/Slemgrim/styx/service"
	"github.com/gorilla/mux"
)

type Upload struct {
	Service service.Attachment
	Store   *gorage.Gorage
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

	err = validateFile(attachment, body)

	if err != nil {
		respondError(w, err)
		return
	}

	savedFile, err := u.Store.Save(attachment.FileName, body, nil)

	if err != nil {
		respondError(w, err)
		return
	}

	u.Service.SetUploadedFile(attachment, savedFile.ID)

}

func respondError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, err)
}

func validateFile(attachment model.Attachment, body []byte) error {
	fileSize := binary.Size(body)
	contentType := http.DetectContentType(body)
	hash := calculateHash(body)

	if fileSize != attachment.Size {
		return errors.New("Filesize doesn't match attachments size")
	}

	if contentType != attachment.MimeType {
		return errors.New("Content type doesn't match attachments content type")
	}

	if hash != attachment.Hash {
		return errors.New("File hash doesn't match attachments hash")
	}

	return nil
}

func calculateHash(body []byte) string {
	h := sha1.New()
	h.Write(body)
	return fmt.Sprintf("%x", h.Sum(nil))
}
