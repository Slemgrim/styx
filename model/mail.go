package model

// Mail defines the mail structure
type Mail struct {
	ID          string   `json:"id"`
	Context     string   `json:"context"`
	Subject     string   `json:"subject"`
	Clients     []Client `json:"clients"`
	Body        Body     `json:"body"`
	Priority    int      `json:"priority"`
	Attachments []Attachment `json:"attachments"`
}

// ClientType defines different available mail headers associated with email addresses
type ClientType string

const (
	CLIENT_TO ClientType = "to"
	CLIENT_FROM ClientType = "from"
	CLIENT_CC ClientType = "cc"
	CLIENT_BCC ClientType = "bcc"
	CLIENT_REPLY_TO ClientType = "reply-to"
	CLIENT_RETURN_PATH ClientType = "return-path"
)

// Client defines the client structure
type Client struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Type  ClientType `json:"type"`
}

type Attachment struct {
	ID           string   `json:"id"`
	OriginalName string `json:"original-name"`
	FileName     string `json:"file-name"`
}

// Body defines the html and plain text fields
type Body struct {
	HTML  string `json:"html"`
	Plain string `json:"plain"`
}

// GetName gets the type identifier of the resource
func (m Mail) GetName() string {
	return "mail"
}

// GetID retrieves the identifier of the mail
func (m Mail) GetID() string {
	return m.ID
}

// SetID sets the identifier of the mail
func (m *Mail) SetID(id string) error {
	m.ID = id
	return nil
}
