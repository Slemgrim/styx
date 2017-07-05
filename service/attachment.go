package service

import (
	"time"

	"github.com/fetzi/styx/model"
	"github.com/fetzi/styx/resource"
	"github.com/google/uuid"
)

/*
	Service for handling attachments
*/
type Attachment struct {
	Resource resource.Attachment
}

/*
	Crete a new attachment
*/
func (a Attachment) Create(attachment model.Attachment) (model.Attachment, error) {
	var err error
	attachment.ID = uuid.New().String()
	attachment.CreatedAt = time.Now()
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
	//Todo
	return a.Resource.Delete(attachment.ID)
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
func (a Attachment) MarkAsUsed(attachment model.Attachment) error {
	time := time.Now()
	attachment.LastUsedAt = &time
	return nil
}
