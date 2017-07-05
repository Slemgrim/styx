package resource

import (
	"errors"
	"fmt"

	"github.com/fetzi/styx/model"

	"github.com/jinzhu/gorm"
)

type Attachment interface {
	Create(attachment model.Attachment) (model.Attachment, error)
	Read(id string) (model.Attachment, error)
	Update(model.Attachment) (model.Attachment, error)
	Delete(id string) error
}

type DbAttachment struct {
	DB *gorm.DB
}

func (a DbAttachment) Init() {
	a.DB.AutoMigrate(&model.Attachment{})
}

func (a DbAttachment) Create(attachment model.Attachment) (model.Attachment, error) {
	result := a.DB.Create(attachment)

	if result.Error != nil {
		return model.Attachment{}, result.Error
	}

	return attachment, nil
}

func (a DbAttachment) Read(id string) (model.Attachment, error) {
	attachment := model.Attachment{}

	notFound := a.DB.Where(model.Attachment{
		ID:        id,
		DeletedAt: nil,
	}, id).First(&attachment).RecordNotFound()

	if notFound {

		return model.Attachment{}, errors.New("attachment not found")
	}

	return attachment, nil

}

func (a DbAttachment) Update(attachment model.Attachment) (model.Attachment, error) {
	result := a.DB.Save(&attachment)

	if result.Error != nil {
		return model.Attachment{}, result.Error
	}

	return attachment, nil

}

func (a DbAttachment) Delete(id string) error {
	fmt.Println("Delete Attachment")
	return nil
}
