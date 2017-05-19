package mailer

import (
	"fmt"
	"github.com/fetzi/styx/config"
	"github.com/fetzi/styx/model"
	"github.com/go-gomail/gomail"
	"errors"
	"os"
)

//Mailer for sending mails
type Mailer struct {
	Dialer *gomail.Dialer
	AttachmentPath string
}

// NewMailer creates a new mailer instance
func NewMailer(smtpConfig config.SMTPConfig, attachmentConfig config.AttachmentConfig) *Mailer {
	dialer := gomail.NewPlainDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.User, smtpConfig.Password)
	return &Mailer{dialer, attachmentConfig.Path}
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
			toList = append(toList, formatEmail(client))
		case model.CLIENT_CC:
			ccList = append(ccList, formatEmail(client))
		case model.CLIENT_BCC:
			bccList = append(bccList, formatEmail(client))
		case model.CLIENT_FROM:
			from = formatEmail(client)
		case model.CLIENT_REPLY_TO:
			replyTo = formatEmail(client)
		case model.CLIENT_RETURN_PATH:
			returnPath = formatEmail(client)
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
		mail.SetHeader("styx-mail-context", data.Context)
	}

	if data.ID == "" {
		return errors.New("Id is missing")
	}

	mail.SetHeader("styx-mail-uuid", data.ID)

	if len(data.Attachments) > 0 {
		for _, attachment := range data.Attachments {
			file := fmt.Sprintf("%s/%s", mailer.AttachmentPath, attachment.FileName)
			if _, err := os.Stat(file); os.IsNotExist(err) {
				return errors.New(fmt.Sprintf("File '%s' doesn't exist", file))
			}
			attachmentIdHeader := map[string][]string{"styx-attachment-uuid": {attachment.ID}}
			mail.Attach(file, gomail.SetHeader(attachmentIdHeader))
		}
	}

	if err := mailer.Dialer.DialAndSend(mail); err != nil {
		return err
	}

	return nil
}

//Format a Client to mail conform string
func formatEmail(client model.Client) string {
	if client.Name == "" {
		return client.Email
	} else {
		return fmt.Sprintf("%s <%s>", client.Name, client.Email)
	}
}
