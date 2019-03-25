package mailgun

import (
	"bytes"
	"context"
	"errors"
	"github.com/mailgun/mailgun-go"
	"html/template"
	"strings"
	"time"
)

type Config struct {
	Html            string
	Text            string
	Domain          string
	PrivateKey      string
	Sender          string
	Subject         string
	EmailValidation mailgun.EmailVerification
}

func (c Config) ParseHtml(data interface{}) error {
	if c.Html != "" {
		return errors.New("set html before using it")
	}
	tmpl, err := template.New("test").Parse(c.Html)
	buf := &bytes.Buffer{}
	if err != nil {
		return err
	}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return err
	}
	c.Html = buf.String()
	return nil
}

func (c Config) Validation(reciver string) (bool, error) {
	v := mailgun.NewEmailValidator(c.PrivateKey)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var email, err = v.ValidateEmail(ctx, reciver, true)
	if err != nil {
		return false, err
	}
	c.EmailValidation = email
	if email.MailboxVerification != "true" {
		return false, nil
	}
	return true, nil
}

func (c Config) Send(reciver ...string) (resp string, id string, err error) {
	var subject = strings.TrimSpace(c.Subject)
	var body = strings.TrimSpace(c.Html)
	if subject == "" {
		return "", "", errors.New("subcjet is empty")
	}
	if body == "" {
		return "", "", errors.New("body html is empty")
	}
	for _, v := range reciver {
		var valid, err = c.Validation(v)
		if !valid || err != nil {
			return "", "", errors.New("invalid e-mail :" + v)
		}
	}
	var mg = mailgun.NewMailgun(c.Domain, c.PrivateKey)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	message := mg.NewMessage(c.Sender, c.Subject, c.Html,reciver[0])
	for i, v := range reciver {
		if i == 0 {continue}
		message.AddBCC(v)
	}
	// Send the message	with a 10 second timeout
	resp, id, err = mg.Send(ctx, message)
	return resp, id, err
}
