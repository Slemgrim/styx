package mailer

import (
	"fmt"
	"github.com/fetzi/styx/config"
	"github.com/fetzi/styx/model"
	"github.com/go-gomail/gomail"
)

//Mailer for sending mails
type Mailer struct {
	Dialer *gomail.Dialer
}

// NewMailer creates a new mailer instance
func NewMailer(config config.SMTPConfig) *Mailer {
	dialer := gomail.NewPlainDialer(config.Host, config.Port, config.User, config.Password)
	return &Mailer{dialer}
}

// Send a mail
func (mailer *Mailer) Send(data model.Mail) error {
	mail := gomail.NewMessage()
	toList := make([]string, 0)
	ccList := make([]string, 0)
	bccList := make([]string, 0)

	for _, client := range data.Clients {
		switch client.Type {
		case model.CLIENT_TO:
			toList = append(toList, formatEmail(client))
		case model.CLIENT_CC:
			ccList = append(ccList, formatEmail(client))
		case model.CLIENT_BCC:
			bccList = append(bccList, formatEmail(client))
		case model.CLIENT_FROM:
			mail.SetHeader("From", formatEmail(client))
		case model.CLIENT_REPLY_TO:
			mail.SetHeader("Reply-To", formatEmail(client))
		case model.CLIENT_RETURN_PATH:
			mail.SetHeader("Return-Path", formatEmail(client))
		}
	}

	mail.SetHeader("To", toList...)
	mail.SetHeader("Cc", ccList...)
	mail.SetHeader("Bcc", bccList...)

	mail.SetHeader("Subject", data.Subject)
	mail.SetBody("text/html", data.Body.HTML)
	mail.AddAlternative("text/plain", data.Body.Plain)

	mail.SetHeader("karriere-mail-context", data.Context)
	mail.SetHeader("karriere-mail-uuid", data.ID)

	if err := mailer.Dialer.DialAndSend(mail); err != nil {
		return err
	}

	return nil
}

//Format a Client to mail conform string
func formatEmail(client model.Client) string {
	return fmt.Sprintf("%s <%s>", client.Name, client.Email)
}
