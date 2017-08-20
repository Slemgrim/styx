package resource

import (
	"errors"
	"fmt"

	"github.com/slemgrim/styx/model"

	"github.com/jinzhu/gorm"
)

type Mail interface {
	Create(mail model.Mail) (model.Mail, error)
	Read(id string) (model.Mail, error)
	Update(model.Mail) (model.Mail, error)
	Delete(id string) error
}

type DbMail struct {
	DB *gorm.DB
}

func (a DbMail) Init() {
	a.DB.AutoMigrate(&model.Mail{})
}

func (a DbMail) Create(mail model.Mail) (model.Mail, error) {
	result := a.DB.Create(mail)

	if result.Error != nil {
		return model.Mail{}, result.Error
	}

	return mail, nil
}

func (a DbMail) Read(id string) (model.Mail, error) {
	mail := model.Mail{}

	notFound := a.DB.Where(model.Mail{
		ID:        id,
		DeletedAt: nil,
	}, id).First(&mail).RecordNotFound()

	if notFound {

		return model.Mail{}, errors.New("mail not found")
	}

	return mail, nil

}

func (a DbMail) Update(mail model.Mail) (model.Mail, error) {
	result := a.DB.Save(&mail)

	if result.Error != nil {
		return model.Mail{}, result.Error
	}

	return mail, nil

}

func (a DbMail) Delete(id string) error {
	fmt.Println("Delete Mail")
	return nil
}
