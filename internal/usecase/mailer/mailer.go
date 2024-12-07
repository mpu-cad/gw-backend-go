package mailer

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"

	"github.com/mpu-cad/gw-backend-go/internal/configs"
	"github.com/mpu-cad/gw-backend-go/internal/models"
)

const (
	// Адрес SMTP-сервера Яндекса
	smtpAuthAddress = "smtp.yandex.ru"
	// Порт и сервер для подключения
	smtpServerAddress = "smtp.yandex.ru:465"
)

type Mailer struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
	channelBroker     chan models.Gmail
}

func New(cfg configs.Mailer, channelBroker chan models.Gmail) *Mailer {
	return &Mailer{
		name:              cfg.Name,
		fromEmailAddress:  cfg.FromEmailAddress,
		fromEmailPassword: cfg.FromEmailPassword,
		channelBroker:     channelBroker,
	}
}

func (m *Mailer) SendEmailToBroker(gmail models.Gmail) error {
	m.channelBroker <- gmail
	return nil
}

func (m *Mailer) SendEmail(gmail models.Gmail) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", m.name, m.fromEmailAddress)
	e.Subject = gmail.Subject
	e.HTML = []byte(gmail.Content)
	e.To = gmail.TO
	e.Cc = gmail.CC
	e.Bcc = gmail.BCC

	// Прикрепление файлов (если есть)
	for i := range gmail.AttachFiles {
		_, err := e.AttachFile(gmail.AttachFiles[i])
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", gmail.AttachFiles[i], err)
		}
	}

	// Используем SMTP с TLS
	smtpAuth := smtp.PlainAuth("", m.fromEmailAddress, m.fromEmailPassword, smtpAuthAddress)
	if err := e.SendWithTLS(smtpServerAddress, smtpAuth, &tls.Config{
		ServerName:         smtpAuthAddress,
		InsecureSkipVerify: false,
	}); err != nil {
		return fmt.Errorf("cannot send email, err: %w", err)
	}
	return nil
}
