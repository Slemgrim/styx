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

// MailStorage tba
type MailStorage struct {
	MailStatusStorage MailStatusStorage
}

// NewMailStorage tba
func NewMailStorage(mailStatusStorage MailStatusStorage) MailStorage {
	return MailStorage{mailStatusStorage}
}

func (s *MailStatusStorage) Insert (m model.MailStatus) string {
	s.db.Create(&m)
	return m.MailID
}
