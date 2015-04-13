package mmr

import (
	"bytes"
	"text/template"
	"net/smtp"
	"strconv"
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
		Subject string
		Body    string
	}
)

const (
	emailTpl = "Content-Type: text/plain; charset=UTF-8\r\nFrom: {{if not .From.Name}}{{.From.Address}}{{else}}{{.From.Name}} <{{.From.Address}}>{{end}}\r\nTo: {{if not .To.Name}}{{.To.Address}}{{else}}{{.To.Name}} <{{.To.Address}}>{{end}}\r\nSubject: {{.Subject}}\r\n\r\n{{.Body}}\r\n"
)

func SendEmail(account *EmailAccount, to *EmailAddress, subject, body string) error {

	tpl, err := template.New("email").Parse(emailTpl)
	if err != nil {
		return err
	}
	
	var doc bytes.Buffer
	err =  tpl.Execute(&doc, &emailData{account.From, to, subject, body})
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", account.Username, account.Password, account.EmailServer)
	return smtp.SendMail(account.EmailServer+":"+strconv.Itoa(account.Port), auth, account.Username, []string{to.Address}, doc.Bytes())
}
