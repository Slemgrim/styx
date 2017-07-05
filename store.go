package styx

import (
	"github.com/Slemgrim/gorage"
	"github.com/Slemgrim/gorage/meta"
	"github.com/Slemgrim/gorage/relation"
	"github.com/Slemgrim/gorage/storage"
	"github.com/fetzi/styx/config"
	"github.com/jinzhu/gorm"
)

func GetBodyStore(config config.FilesConfig, db *gorm.DB) *gorage.Gorage {
	s := storage.Io{
		BasePath:   config.BodyPath,
		DirLength:  6,
		BufferSize: 1024,
	}

	r := relation.NewDb("mail_body_rel", db)
	m := meta.NewDb("mail_body", db)

	gorage := gorage.NewGorage(s, r, m)

	return gorage
}

func GetAttachmentStore(config config.FilesConfig, db *gorm.DB) *gorage.Gorage {
	s := storage.Io{
		BasePath:   config.AttachmentPath,
		DirLength:  6,
		BufferSize: 1024,
	}

	r := relation.NewDb("attachment_store_rel", db)
	m := meta.NewDb("attachment_store", db)

	gorage := gorage.NewGorage(s, r, m)

	return gorage
}
