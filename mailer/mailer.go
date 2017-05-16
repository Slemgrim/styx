package mailer

import (
	"github.com/fetzi/styx/model"
	"fmt"
	"github.com/domodwyer/mailyak"
	"github.com/fetzi/styx/config"
	"net/smtp"
)

//Send mail to all defined clients
func SendMail(data model.Mail, config config.SMTPConfig) error {
	mail := mailyak.New(fmt.Sprintf("%s:%d", config.Host, config.Port), smtp.PlainAuth(config.Identity, config.User, config.Password, config.Host))

	toList := make([]string, 0)
	bccList := make([]string, 0)

	for _, client := range data.Clients {
		switch client.Type {
		case model.CLIENT_TO:
			toList = append(toList, formatEmail(client))
		case model.CLIENT_BC:
		//TODO how to send bc with mailyak
		case model.CLIENT_BCC:
			bccList = append(toList, formatEmail(client))
		case model.CLIENT_FROM:
			mail.From(client.Email)
			mail.FromName(client.Name)
		case model.CLIENT_REPLY_TO:
			mail.ReplyTo(formatEmail(client))
		case model.CLIENT_RETURN_PATH:
		//TODO how to send return path with mailyak
		}
	}

	mail.To(toList...)
	mail.Bcc(bccList...)
	mail.Subject(data.Subject)
	mail.Plain().Set(data.Body.Plain)
	mail.HTML().Set(data.Body.HTML)

	if err := mail.Send(); err != nil {
		return err
	}

	return nil
}

func formatEmail(client model.Client) string {
	return fmt.Sprintf("%s <%s>", client.Name, client.Email)
}