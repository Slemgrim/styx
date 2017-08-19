package model

import (
	"testing"

	"gopkg.in/go-playground/validator.v9"
	"os"
)

var v *validator.Validate = validator.New()

func TestMain(m *testing.M){
	v.RegisterStructValidation(ValidateBody, Body{})
	v.RegisterStructValidation(ValidateAddress, Address{})

	os.Exit(m.Run())
}

func TestMinimalMail(t *testing.T) {
	address := Address{
		Name: "Rick Sanchez",
		Mail: "RickSanchez@example.com",
	}

	to := []Address{}
	to = append(to, address)

	body := Body{
		HTML: "<h1>Eat lasers!<h1>",
		Plain: "Eat Lasers!",
	}

	mail := Mail{
		Subject: "Eat Lasers!",
		To: to,
		From: address,
		Body: body,
	}

	err := v.Struct(mail)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestBodyIsRequired(t *testing.T) {
	address := Address{
		Name: "Rick Sanchez",
		Mail: "RickSanchez@example.com",
	}

	to := []Address{}
	to = append(to, address)

	mail := Mail{
		Subject: "Eat Lasers!",
		To: to,
		From: address,
	}

	err := v.Struct(mail)
	if err == nil {
		t.Fatal(err.Error())
	}
}

func TestSubjectIsRequired(t *testing.T) {
	address := Address{
		Name: "Rick Sanchez",
		Mail: "RickSanchez@example.com",
	}

	to := []Address{}
	to = append(to, address)

	body := Body{
		HTML: "<h1>Eat lasers!<h1>",
		Plain: "Eat Lasers!",
	}

	mail := Mail{
		To: to,
		From: address,
		Body: body,
	}

	err := v.Struct(mail)
	if err == nil {
		t.Fatal(err.Error())
	}
}


func TestAtLeastOnToIsRequired(t *testing.T) {
	address := Address{
		Name: "Rick Sanchez",
		Mail: "RickSanchez@example.com",
	}

	body := Body{
		HTML: "<h1>Eat lasers!<h1>",
		Plain: "Eat Lasers!",
	}

	mail := Mail{
		Subject: "Eat Lasers!",
		From: address,
		Body: body,
	}

	err := v.Struct(mail)
	if err == nil {
		t.Fatal(err.Error())
	}
}

func TestFromIsRequired(t *testing.T) {
	//To do
}