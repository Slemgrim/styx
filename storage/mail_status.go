package storage

import (
	"fmt"

	"github.com/fetzi/styx/model"
	"github.com/jinzhu/gorm"
)

// MailStatusStorage tba
type MailStatusStorage struct {
	db *gorm.DB
}

// NewMailStatusStorage tba
func NewMailStatusStorage(db *gorm.DB) MailStatusStorage {
	table := &model.MailStatus{}

	if !db.HasTable(table) {
		fmt.Println("creating table MailStatus")

		db.CreateTable(table)
	}

	return MailStatusStorage{db}
}

// GetOne tba
func (s MailStatusStorage) GetOne(id string) (model.MailStatus, error) {
	mailStatus := model.MailStatus{}
	notFound := s.db.Where(model.MailStatus{
		MailID: id,
	}).First(&mailStatus).RecordNotFound()

	if notFound {
		return mailStatus, fmt.Errorf("Mail Status with id %s not found", id)
	}

	return mailStatus, nil

}
