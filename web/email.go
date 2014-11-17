package mmr

import (
	"bytes"
	"net/smtp"
	"strconv"
	"html/template"
)

type (
	EmailAccount struct {
		EmailServer string
		Port        int
		Username    string
		Password    string
		From        string
	}

	emailData struct {
		From    string
		To      string
		Subject string
		Body    string
	}
)

func SendEmail(account *EmailAccount, tpl *template.Template, to, subject, body string) error {

	var doc bytes.Buffer
	err := tpl.Execute(&doc, &emailData{account.From, to, subject, body})
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", account.Username, account.Password, account.EmailServer)
	return smtp.SendMail(account.EmailServer+":"+strconv.Itoa(account.Port), auth, account.Username, []string{to}, doc.Bytes())
}
