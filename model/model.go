package model

// Mail defines the mail structure
type Mail struct {
	Context     string
	Subject     string
	Clients     []Client
	Body        Body
	Priority    int
	Attachments []string
}

// Client defines the client structure
type Client struct {
	Name  string
	Email string
	Type  string
}

// Body defines the html and plain text fields
type Body struct {
	HTML  string
	Plain string
}

type MailStatus struct {
	MailID  string `json:"-" gorm:"primary_key"`
	Subject string `json:"subject"`
	From    string `json:"from"`
	To      string `json:"to"`
	Created int64  `json:"created_at"`
	Sent    int64  `json:"sent_at"`
}
