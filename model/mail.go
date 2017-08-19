package model

import (
	"time"
	"gopkg.in/go-playground/validator.v9"
	"github.com/badoux/checkmail"
)

type Mail struct {
	ID         string `jsonapi:"primary,attachments" gorm:"primary_key"`
	Subject    string `jsonapi:"attr,subject" validate:"required"`
	Body	   Body `jsonapi:"attr,body" validate:"required,dive"`
	Attachments []Attachment

	To			[]Address `jsonapi:"attr,to" validate:"required,dive,required,gte=1"`
	Cc			[]Address `jsonapi:"attr,cc" validate:"dive"`
	Bcc			[]Address `jsonapi:"attr,bcc" validate:"dive"`

	From		Address `jsonapi:"attr,from" validate:"required"`
	ReplyTo		Address `jsonapi:"attr,reply-to" validate:"omitempty"`
	ReturnPath	Address `jsonapi:"attr,return-path" validate:"omitempty"`

	Headers		[]Header `jsonapi:"attr,headers" validate:"dive"`

	CreatedAt  time.Time
	SentAt *time.Time
	DeletedAt  *time.Time
}

type Body struct {
	Plain string `json:"plain"`
	HTML string `json:"html" validate:""`
}

type Address struct {
	Name string `json:"name"`
	Mail string `json:"mail"`
}

type Header struct {
	Name string `json:"name" validate:"required"`
	Value []string `json:"value"`
}


func ValidateBody(sl validator.StructLevel) {
	body := sl.Current().Interface().(Body)
	if len(body.HTML) == 0 && len(body.Plain) == 0 {
		sl.ReportError(body.HTML, "HTML", "html", "htmlorplain", "")
		sl.ReportError(body.Plain, "Plain", "plain", "htmlorplain", "")
	}
}

func ValidateAddress(sl validator.StructLevel) {
	address := sl.Current().Interface().(Address)

	if address.Mail == "" {
		return
	}

	if err := checkmail.ValidateFormat(address.Mail); err != nil {
		sl.ReportError(address.Mail, "Mail", "mail", "novalidemail", "")
	}
}
