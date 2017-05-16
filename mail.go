package main

import (
	"github.com/fetzi/styx/model"
	"github.com/fetzi/styx/mailer"
	"github.com/fetzi/styx/config"
	"log"
)

//getTestMail creates an example mal
func getTestMail() model.Mail {
	mail := model.Mail{
		ID: "test-mail-1234-1234",
		Subject: "Test Mail",
		Context: "TEST_MAIL",
		Body: model.Body{
			HTML: "<h1>Test Mail</h1>",
			Plain: "Test Mail",
		},
		Clients: make([]model.Client, 0),
	}

	to1 := model.Client{
		Type: model.CLIENT_TO,
		Name: "Markus Ritberger",
		Email: "markus.ritberger@karriere.at",
	}

	to2 := model.Client{
		Type: model.CLIENT_TO,
		Name: "Markus Ritberger",
		Email: "markus.ritberger@slemgrim.com",
	}

	from := model.Client{
		Type: model.CLIENT_FROM,
		Name: "Styx der hund",
		Email: "styx@karriere.at",
	}

	bcc1 := model.Client{
		Type: model.CLIENT_BCC,
		Name: "Styx der hund",
		Email: "styx@karriere.at",
	}

	bcc2 := model.Client{
		Type: model.CLIENT_BCC,
		Name: "Styx der hund",
		Email: "styx@karriere.at",
	}

	replyTo := model.Client{
		Type: model.CLIENT_REPLY_TO,
		Name: "Styx der hund",
		Email: "styx@karriere.at",
	}

	returnPath := model.Client{
		Type: model.CLIENT_RETURN_PATH,
		Name: "Styx der hund",
		Email: "styx@karriere.at",
	}

	mail.Clients = append(mail.Clients, to1)
	mail.Clients = append(mail.Clients, to2)
	mail.Clients = append(mail.Clients, from)
	mail.Clients = append(mail.Clients, bcc1)
	mail.Clients = append(mail.Clients, bcc2)
	mail.Clients = append(mail.Clients, replyTo)
	mail.Clients = append(mail.Clients, returnPath)

	return mail;
}

func main (){

	config, err := config.ReadConfig("config.json")

	if err != nil {
		log.Fatal(err)
	}

	mail := getTestMail()
	mailer.SendMail(mail, config.SMTP)
}