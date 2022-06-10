package main

import (
	"fmt"
	"time"

	"gopkg.in/gomail.v2"
)

type email struct {
	user     string
	password string
}

func newEmail(user, password string) *email {
	return &email{
		user:     user,
		password: password,
	}
}

func (e *email) send(to []string, msg string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.user)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", fmt.Sprintf("%s's status - %s", name, time.Now().Format("2006-01-02")))
	m.SetBody("text/plain", msg)

	d := gomail.NewDialer("smtp.gmail.com", 587, e.user, e.password)
	return d.DialAndSend(m)
}
