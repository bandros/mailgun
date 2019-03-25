package mailgun

import (
	"context"
	"github.com/mailgun/mailgun-go"
	"time"
)

type Config struct {
	Html string
	Text string
	Domain string
	PrivateKey string
	Sender string
	Subject string
}

func (c Config) Validation(reciver string) (bool,error) {
	v := mailgun.NewEmailValidator(c.PrivateKey)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var email, err = v.ValidateEmail(ctx, reciver, false)
	if err != nil {
		return false,err
	}
	return email.IsValid,nil
}

func (c Config) Send(reciver ...string) (resp string,id string, err error){
	var mg = mailgun.NewMailgun(c.Domain, c.PrivateKey)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	message := mg.NewMessage(c.Sender, c.Subject, c.Html, reciver...)
	// Send the message	with a 10 second timeout
	resp,id,err = mg.Send(ctx, message)
	return resp,id,err
}