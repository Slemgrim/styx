package mailer

import (
	"fmt"
	"github.com/fetzi/styx/config"
	"github.com/fetzi/styx/model"
	"github.com/go-gomail/gomail"
	"errors"
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
	var from string;
	var replyTo string;
	var returnPath string;

	for _, client := range data.Clients {
		switch client.Type {
		case model.CLIENT_TO:
			email, err := formatEmail(client)
			if err == nil {
				toList = append(toList, email)
			}
		case model.CLIENT_CC:
			email, err := formatEmail(client)
			if err == nil {
				ccList = append(ccList, email)
			}
		case model.CLIENT_BCC:
			email, err := formatEmail(client)
			if err == nil {
				bccList = append(bccList, email)
			}
		case model.CLIENT_FROM:
			email, err := formatEmail(client)
			if err == nil {
				from = email
			}
		case model.CLIENT_REPLY_TO:
			email, err := formatEmail(client)
			if err == nil {
				replyTo = email
			}
		case model.CLIENT_RETURN_PATH:
			email, err := formatEmail(client)
			if err == nil {
				returnPath = email
			}
		}
	}

	if len(toList) == 0 {
		return errors.New("To header missing")
	}
	mail.SetHeader("To", toList...)

	if from == "" {
		return errors.New("From header missing")
	}

	mail.SetHeader("From", from)

	if len(ccList) > 0 {
		mail.SetHeader("Cc", ccList...)
	}

	if len(bccList) > 0 {
		mail.SetHeader("Bcc", bccList...)
	}

	if replyTo != "" {
		mail.SetHeader("Reply-To", replyTo)
	}

	if returnPath != "" {
		mail.SetHeader("Return-Path", returnPath)
	}

	if data.Subject == "" {
		return errors.New("Subject is missing")
	}
	mail.SetHeader("Subject", data.Subject)

	if data.Body.HTML == "" && data.Body.Plain == "" {
		return errors.New("No body was provided")
	}

	if data.Body.HTML != "" {
		mail.SetBody("text/html", data.Body.HTML)
	}

	if data.Body.Plain != "" {
		mail.AddAlternative("text/plain", data.Body.Plain)
	}

	if data.Context != "" {
		mail.SetHeader("karriere-mail-context", data.Context)
	}

	if data.ID == "" {
		return errors.New("Id is missing")
	}

	mail.SetHeader("karriere-mail-uuid", data.ID)

	if err := mailer.Dialer.DialAndSend(mail); err != nil {
		return err
	}

	return nil
}

//Format a Client to mail conform string
func formatEmail(client model.Client) (string, error) {
	if client.Email == "" {
		return "", errors.New("Missing email")
	}
	if client.Name == "" {
		return client.Email, nil
	} else {
		return fmt.Sprintf("%s <%s>", client.Name, client.Email), nil
	}
}
