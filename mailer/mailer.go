package mailer

import (
	"bytes"
	"errors"
	"github.com/Slemgrim/gorage"
	"github.com/Slemgrim/styx/config"
	"github.com/Slemgrim/styx/model"
	"github.com/go-gomail/gomail"
	"io"
	"time"
)

type Mailer struct {
	Dialer *gomail.Dialer
	Store  *gorage.Gorage
}

// Creates a new mailer instance
func New(smtpConfig config.SMTPConfig, store *gorage.Gorage) *Mailer {
	dialer := gomail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.User, smtpConfig.Password)
	return &Mailer{dialer, store}
}

// Send a mail
func (mailer *Mailer) Send(mail model.Mail) error {
	message := gomail.NewMessage()

	to, err := getAddressList(mail.To, message)
	if err != nil {
		return err
	}

	cc, err := getAddressList(mail.Cc, message)
	if err != nil {
		return err
	}

	bcc, err := getAddressList(mail.Bcc, message)
	if err != nil {
		return err
	}

	from, _ := getAddress(mail.From, message)
	replyTo, _ := getAddress(mail.ReplyTo, message)
	returnPath, _ := getAddress(mail.ReturnPath, message)

	message.SetHeader("To", to...)
	message.SetHeader("From", from)

	message.SetHeader("Cc", cc...)
	message.SetHeader("Bcc", bcc...)

	if replyTo != "" {
		message.SetHeader("Reply-To", replyTo)
	}

	if returnPath != "" {
		message.SetHeader("Return-Path", returnPath)
	}

	message.SetHeader("Subject", mail.Subject)

	if mail.Body.HTML != "" {
		message.SetBody("text/html", mail.Body.HTML)
	}

	if mail.Body.Plain != "" {
		message.AddAlternative("text/plain", mail.Body.Plain)
	}

	message.SetHeader("styx-mail-uuid", mail.ID)
	message.SetHeader("styx-mail-date", message.FormatDate(time.Now()))

	mailer.addAttachments(message, mail.Attachments, mailer)

	if err := mailer.Dialer.DialAndSend(message); err != nil {
		return err
	}

	return nil
}

func getAddressList(addressList []model.Address, message *gomail.Message) ([]string, error) {
	list := make([]string, 0)

	for _, address := range addressList {
		email, err := getAddress(address, message)
		if err == nil {
			list = append(list, email)
		}
	}

	return list, nil
}

func getAddress(address model.Address, message *gomail.Message) (string, error) {
	if address.Address == "" {
		return "", errors.New("Missing email")
	}

	return message.FormatAddress(address.Address, address.Name), nil
}

func (m *Mailer) addAttachments(mail *gomail.Message, attachments []*model.Attachment, mailer *Mailer) error {
	if len(attachments) > 0 {
		for _, attachment := range attachments {
			attachmentIdHeader := map[string][]string{"styx-attachment-uuid": {attachment.ID}}
			f, err := m.Store.Load(attachment.FileId)

			if err != nil {
				return err
			}

			r := bytes.NewBuffer(f.Content)
			mail.Attach(attachment.FileName, gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := io.Copy(w, r)
				return err
			}), gomail.Rename(attachment.FileName), gomail.SetHeader(attachmentIdHeader))
		}
	}
	return nil
}
