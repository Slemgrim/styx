package mailer

import (
	"errors"
	"github.com/Slemgrim/styx/config"
	"github.com/Slemgrim/styx/model"
	"github.com/go-gomail/gomail"
	"time"
)

type Mailer struct {
	Dialer *gomail.Dialer
}

// Creates a new mailer instance
func New(smtpConfig config.SMTPConfig) *Mailer {
	dialer := gomail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.User, smtpConfig.Password)
	return &Mailer{dialer}
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

	//addAttachments(message, mail.Attachments, mailer)

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

/*
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
*/
