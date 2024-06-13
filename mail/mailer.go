package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/go-gomail/gomail"
	"html/template"
	"log"
	"time"
)

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

type BreadcrumbReport struct {
	Header      string
	Report      []string
	Attachments []string
}

func (b *BreadcrumbReport) Append(line string) {
	b.Report = append(b.Report, line)
}

func NewRequest(to []string, subject string, from string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		from:    from,
	}
}

func Mail(config *Config, report *BreadcrumbReport) error {
	r := NewRequest(config.Message.To, config.Message.Subject, config.Email.Email)
	return r.Send(config, report)
}

func (r *Request) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func (r *Request) sendMail(config *Config, attachments []string) error {
	d, err := GetDialer(config)
	if err != nil {
		return err
	}
	m := gomail.NewMessage()
	m.SetAddressHeader("From", r.from, r.from)
	m.SetHeader("To", r.to[0])
	m.SetDateHeader("X-Date", time.Now())
	m.SetHeader("Subject", r.subject)
	m.SetBody("text/html", r.body)
	for _, attachment := range attachments {
		m.Attach(attachment)
	}
	return d.DialAndSend(m)
}

func (r *Request) Send(config *Config, items interface{}) error {
	err := r.parseTemplate(config.Message.Template, items)
	if err != nil {
		return fmt.Errorf("could not parse HTML template %v", err)
	}

	attachments := make([]string, 0)
	report, ok := items.(*BreadcrumbReport)
	if ok {
		attachments = report.Attachments
	}

	err = r.sendMail(config, attachments)
	if err != nil {
		return fmt.Errorf("Failed to send the email to %s\n%v\n", r.to, err)
	}
	log.Printf("Email has been sent to %s\n", r.to)
	return nil
}

func GetDialer(cfg *Config) (*gomail.Dialer, error) {
	email := cfg.Email
	if email.Server == "" {
		err := fmt.Errorf("SMTP server not configured")
		return nil, err
	}

	d := gomail.NewDialer(email.Server, email.Port, email.Email, email.Password)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	/*
		if cfg.AuthType == config.EmailAuthCramMD5 {
			d.Auth = smtp.CRAMMD5Auth(cfg.User, cfg.Password)
		}
	*/
	return d, nil
}
