package utils

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
)

type EmailTransporter struct {
	Host        string
	Port        int
	User        string
	Password    string
	From        string
	FrontendURL string
}

// Load email transport from ENV (only once)
func NewEmailTransporter() (*EmailTransporter, error) {
	host := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	portStr := strings.TrimSpace(os.Getenv("SMTP_PORT"))
	user := strings.TrimSpace(os.Getenv("SMTP_USER"))
	pass := strings.TrimSpace(os.Getenv("SMTP_PASSWORD"))
	from := strings.TrimSpace(os.Getenv("EMAIL_FROM"))
	frontend := strings.TrimSpace(os.Getenv("FRONTEND_URL"))

	if host == "" || portStr == "" || user == "" || pass == "" || from == "" {
		return nil, errors.New("missing SMTP configuration environment variables")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.New("SMTP_PORT must be a valid number")
	}

	return &EmailTransporter{
		Host:        host,
		Port:        port,
		User:        user,
		Password:    pass,
		From:        from,
		FrontendURL: frontend,
	}, nil
}

// Sends a plain-text email
func (t *EmailTransporter) SendPlain(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", t.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(t.Host, t.Port, t.User, t.Password)

	return d.DialAndSend(m)
}

// Sends an HTML email
func (t *EmailTransporter) SendHTML(to, subject, html string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", t.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", html)

	d := gomail.NewDialer(t.Host, t.Port, t.User, t.Password)

	return d.DialAndSend(m)
}
