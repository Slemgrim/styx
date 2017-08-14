package storage

import "github.com/slemgrim/styx/model"

// MailStorage tba
type MailStorage struct {
	MailStatusStorage MailStatusStorage
}

// NewMailStorage creates a new MailStorage instance
func NewMailStorage(mailStatusStorage MailStatusStorage) MailStorage {
	return MailStorage{mailStatusStorage}
}

// Insert inserts the given mail status object into the database
func (s *MailStatusStorage) Insert(m model.MailStatus) string {
	s.db.Create(&m)
	return m.MailID
}
