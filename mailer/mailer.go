package mailer

import (
	"fmt"
	"github.com/slemgrim/styx/config"
	"github.com/slemgrim/styx/model"
	"github.com/go-gomail/gomail"
	"errors"
	"os"
	"time"
)

//Mailer for sending mails
type Mailer struct {
	Dialer         *gomail.Dialer
	AttachmentPath string
}

// NewMailer creates a new mailer instance
func NewMailer(smtpConfig config.SMTPConfig, attachmentConfig config.AttachmentConfig) *Mailer {
	dialer := gomail.NewPlainDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.User, smtpConfig.Password)
	return &Mailer{dialer, attachmentConfig.Path}
}

// Send a mail
func (mailer *Mailer) Send(mail model.Mail) error {
	message := gomail.NewMessage()
	toList := make([]string, 0)
	ccList := make([]string, 0)
	bccList := make([]string, 0)
	var from string;
	var replyTo string;
	var returnPath string;

	for _, client := range mail.Clients {
		switch client.Type {
		case model.CLIENT_TO:
			toList = addClientToList(toList, client, message)
		case model.CLIENT_CC:
			ccList = addClientToList(toList, client, message)
		case model.CLIENT_BCC:
			bccList = addClientToList(toList, client, message)
		case model.CLIENT_FROM:
			from = setValidClient(client, message)
		case model.CLIENT_REPLY_TO:
			replyTo = setValidClient(client, message)
		case model.CLIENT_RETURN_PATH:
			returnPath = setValidClient(client, message)
		}
	}

	if len(toList) == 0 &&  len(ccList) == 0 && len(bccList) == 0{
		return errors.New("A mail needs at least on to, cc or bcc email")
	}
	message.SetHeader("To", toList...)

	if from == "" {
		return errors.New("From header missing")
	}

	message.SetHeader("From", from)

	if len(ccList) > 0 {
		message.SetHeader("Cc", ccList...)
	}

	if len(bccList) > 0 {
		message.SetHeader("Bcc", bccList...)
	}

	if replyTo != "" {
		message.SetHeader("Reply-To", replyTo)
	}

	if returnPath != "" {
		message.SetHeader("Return-Path", returnPath)
	}

	if mail.Subject == "" {
		return errors.New("Subject is missing")
	}
	message.SetHeader("Subject", mail.Subject)

	if mail.Body.HTML == "" && mail.Body.Plain == "" {
		return errors.New("No body was provided")
	}

	if mail.Body.HTML != "" {
		message.SetBody("text/html", mail.Body.HTML)
	}

	if mail.Body.Plain != "" {
		message.AddAlternative("text/plain", mail.Body.Plain)
	}

	if mail.Context != "" {
		message.SetHeader("styx-mail-context", mail.Context)
	}

	if mail.ID == "" {
		return errors.New("Id is missing")
	}

	message.SetHeader("styx-mail-uuid", mail.ID)
	message.SetHeader("styx-mail-date", message.FormatDate(time.Now()))

	addAttachments(message, mail.Attachments, mailer)

	if err := mailer.Dialer.DialAndSend(message); err != nil {
		return err
	}

	return nil
}


//Format a Client to mail conform string
func formatEmail(client model.Client, message *gomail.Message) (string, error) {
	if client.Email == "" {
		return "", errors.New("Missing email")
	}
	return message.FormatAddress(client.Email, client.Name), nil
}

func addClientToList(list []string, client model.Client, message *gomail.Message) []string {
	email, err := formatEmail(client, message)
	if err == nil {
		list = append(list, email)
	}

	return list
}

func setValidClient(client model.Client, message *gomail.Message) string {
	email, err := formatEmail(client, message)
	if err == nil {
		return email
	}

	return ""
}


func addAttachments(mail *gomail.Message, attachments []model.Attachment, mailer *Mailer) error {
	if len(attachments) > 0 {
		for _, attachment := range attachments {
			file := fmt.Sprintf("%s/%s", mailer.AttachmentPath, attachment.FileName)
			if _, err := os.Stat(file); os.IsNotExist(err) {
				return errors.New(fmt.Sprintf("File '%s' doesn't exist", file))
			}
			attachmentIdHeader := map[string][]string{"styx-attachment-uuid": {attachment.ID}}
			mail.Attach(file, gomail.Rename(attachment.OriginalName), gomail.SetHeader(attachmentIdHeader))
		}
	}
	return nil
}