package service

import (
	"time"

	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Slemgrim/gorage"
	"github.com/Slemgrim/styx/model"
	"github.com/Slemgrim/styx/resource"
	"github.com/google/uuid"
	"net/http"
)

/*
	Service for handling attachments
*/
type Attachment struct {
	Resource resource.Attachment
	Store    *gorage.Gorage
}

/*
	Crete a new attachment
*/
func (a Attachment) Create(attachment model.Attachment) (model.Attachment, error) {
	var err error
	attachment.ID = uuid.New().String()
	attachment.CreatedAt = time.Now()
	attachment.DeletedAt = time.Time{}
	attachment.LastUsedAt = time.Time{}
	attachment.IsUploaded = false

	attachment, err = a.Resource.Create(attachment)

	return attachment, err
}

/*
	Load an attachment by its id
*/
func (a Attachment) Load(id string) (model.Attachment, error) {

	attachment, err := a.Resource.Read(id)

	if err != nil {
		return model.Attachment{}, err
	}

	return attachment, nil
}

/*
	Delete an attachment
*/
func (a Attachment) Delete(attachment model.Attachment) error {
	attachment.DeletedAt = time.Now()
	_, err := a.Resource.Update(attachment)
	return err
}

/*
	Add id of uploaded file to attachment
*/
func (a Attachment) SetUploadedFile(attachment model.Attachment, fileId string) (model.Attachment, error) {
	attachment.FileId = fileId
	attachment.IsUploaded = true
	return a.Resource.Update(attachment)
}

/*
	Mark an attachment as used so it can be deleted after some time
*/
func (a Attachment) MarkAsUsed(attachment model.Attachment) (model.Attachment, error) {
	attachment.LastUsedAt = time.Now()
	return a.Resource.Update(attachment)
}

func (a Attachment) ValidateFile(attachment model.Attachment, body []byte) error {
	fileSize := binary.Size(body)
	contentType := http.DetectContentType(body)
	hash := a.CalculateHash(body)

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

func (a Attachment) CalculateHash(body []byte) string {
	h := sha1.New()
	h.Write(body)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (a Attachment) Upload(attachment model.Attachment, body []byte, context interface{}) error {

	savedFile, err := a.Store.Save(attachment.FileName, body, nil)

	if err != nil {
		return err
	}

	_, err = a.SetUploadedFile(attachment, savedFile.ID)

	return err
}
