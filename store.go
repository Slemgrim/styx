package styx

import (
	"github.com/Slemgrim/gorage"
	"github.com/Slemgrim/gorage/meta"
	"github.com/Slemgrim/gorage/relation"
	"github.com/Slemgrim/gorage/storage"
	"github.com/Slemgrim/styx/config"
	"gopkg.in/mgo.v2"
)

func GetAttachmentStore(config config.FilesConfig, database *mgo.Database) *gorage.Gorage {
	s := storage.Io{
		BasePath:   config.AttachmentPath,
		DirLength:  6,
		BufferSize: 1024,
	}

	r := relation.Mongo{Collection: database.C("attachment_relation")}
	m := meta.Mongo{Collection: database.C("attachment_meta")}

	gorage := gorage.NewGorage(s, r, m)

	return gorage
}
