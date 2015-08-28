package mmr

import (
	"gopkg.in/gomail.v1"
	"strings"
	"errors"
)

type (
	EmailAddress struct {
		Name    string
		Address string
	}

	EmailAccount struct {
		EmailServer string
		Port        int
		Username    string
		Password    string
		From        *EmailAddress
	}

	emailData struct {
		From    *EmailAddress
		To      *EmailAddress
		ReplyTo *EmailAddress
		Subject string
		Body    string
	}
)

func SendEmail(account *EmailAccount, to, replyTo *EmailAddress, subject, contentType, body string) error {

	if account.From == nil {
		return errors.New("You need to specify a From address.")
	}
	
	if to == nil {
		return errors.New("You need to specify a To address.")
	}

	msg := gomail.NewMessage()
	if strings.Contains(account.From.Name, ",") {
		msg.SetHeader("From", account.From.Address)
	} else {
		msg.SetAddressHeader("From", account.From.Address, account.From.Name)
	}
	if strings.Contains(to.Name, ",") {
		msg.SetHeader("To", to.Address)
	} else {
		msg.SetAddressHeader("To", to.Address, to.Name)
	}
	if replyTo != nil {
		if strings.Contains(replyTo.Name, ",") {
			msg.SetHeader("Reply-To", replyTo.Address)
		} else {
			msg.SetAddressHeader("Reply-To", replyTo.Address, replyTo.Name)
		}
	}
	msg.SetHeader("Subject", subject)
	msg.SetBody(contentType, body)

	mailer := gomail.NewMailer(account.EmailServer, account.Username, account.Password, account.Port)
	return mailer.Send(msg)
}
