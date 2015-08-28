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

func quote(name string) string {

	return "\"" + strings.Replace(name, "\"", "'", -1) + "\""
} 

func SendEmail(account *EmailAccount, to, replyTo *EmailAddress, subject, contentType, body string) error {

	if account.From == nil {
		return errors.New("You need to specify a From address.")
	}
	
	if to == nil {
		return errors.New("You need to specify a To address.")
	}

	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", account.From.Address, quote(account.From.Name))
	msg.SetAddressHeader("To", to.Address, quote(to.Name))
	if replyTo != nil {
		msg.SetAddressHeader("Reply-To", replyTo.Address, quote(replyTo.Name))
	}
	msg.SetHeader("Subject", subject)
	msg.SetBody(contentType, body)

	mailer := gomail.NewMailer(account.EmailServer, account.Username, account.Password, account.Port)
	return mailer.Send(msg)
}
