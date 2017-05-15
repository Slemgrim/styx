package model

// Mail defines the mail structure
type Mail struct {
	ID          string   `json:"-"`
	Context     string   `json:"context"`
	Subject     string   `json:"subject"`
	Clients     []Client `json:"clients"`
	Body        Body     `json:"body"`
	Priority    int      `json:"priority"`
	Attachments []string `json:"attachments"`
}

// Client defines the client structure
type Client struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Type  string `json:"type"`
}

// Body defines the html and plain text fields
type Body struct {
	HTML  string `json:"html"`
	Plain string `json:"plain"`
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
