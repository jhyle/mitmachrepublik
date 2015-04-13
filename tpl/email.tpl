Content-Type: text/plain; charset=UTF-8
From: {{if not .From.Name}}{{.From.Address}}{{else}}{{.From.Name}} <{{.From.Address}}>{{end}}
To: {{if not .To.Name}}{{.To.Address}}{{else}}{{.To.Name}} <{{.To.Address}}>{{end}}
Subject: {{.Subject}}

{{.Body}}
