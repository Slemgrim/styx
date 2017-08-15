package model

import (
	"time"
)

type Mail struct {
	ID         string `jsonapi:"primary,attachments" gorm:"primary_key"`
	Subject    string `jsonapi:"attr,subject" validate:"required"`
	Body	   Body `jsonapi:"attr,body" validate:"required"`
	Attachments []Attachment  `jsonapi:"attr,attachments"`

	To			[]Address `jsonapi:"attr,to" validate:"required"`
	Cc			[]Address `jsonapi:"attr,cc"`
	Bcc			[]Address `jsonapi:"attr,bcc"`

	From		Address `jsonapi:"attr,from" validate:"required"`
	ReplyTo		Address `jsonapi:"attr,reply-to"`
	ReturnPath	Address `jsonapi:"attr,return-path"`

	Headers		[]Header `jsonapi:"attr,headers"`

	CreatedAt  time.Time
	SentAt *time.Time
	DeletedAt  *time.Time
}

type Body struct {
	Plain string `json:"plain"`
	HTML string `json:"html"`
}

type Address struct {
	Name string `json:"name"`
	Mail string `json:"mail" validate:"required"`
}

type Header struct {
	Name string `json:"name" validate:"required"`
	Value []string `json:"value"`
}