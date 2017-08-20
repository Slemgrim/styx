package model

import (
	"fmt"
	"time"

	"github.com/Slemgrim/jsonapi"
)

type Attachment struct {
	ID         string `jsonapi:"primary,attachments" gorm:"primary_key"`
	FileId     string
	Size       int    `jsonapi:"attr,size" validate:"required"`
	FileName   string `jsonapi:"attr,file-name" validate:"required"`
	MimeType   string `jsonapi:"attr,mime-type" validate:"required"`
	Hash       string `jsonapi:"attr,hash"`
	CreatedAt  time.Time
	LastUsedAt *time.Time
	DeletedAt  *time.Time
	IsUploaded bool `jsonapi:"attr,is-uploaded" `
}

func (a Attachment) JSONAPILinks() *jsonapi.Links {

	if !a.IsUploaded {
		return &jsonapi.Links{
			"upload": fmt.Sprintf("/upload/%s", a.ID),
		}
	}
	return &jsonapi.Links{}
}
