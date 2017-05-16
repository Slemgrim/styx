package model

// MailStatus defines the status entry of a mail
type MailStatus struct {
	MailID  string `json:"-" gorm:"primary_key"`
	Subject string `json:"subject"`
	From    string `json:"from"`
	To      string `json:"to"`
	Created int64  `json:"created_at"`
	Sent    int64  `json:"sent_at"`
}

// GetName gets the type identifier of the resource
func (m MailStatus) GetName() string {
	return "mail-status"
}

// GetID retrieves the identifier of the mail
func (m MailStatus) GetID() string {
	return m.MailID
}

// SetID sets the identifier of the mail
func (m *MailStatus) SetID(id string) error {
	m.MailID = id
	return nil
}
