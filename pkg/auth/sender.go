package auth

import (
	"context"
	"fmt"
	"gopkg.in/gomail.v2"
)

const (
	subject = "Ваш код для входа"
)

type DefaultSender struct {
	config *SenderConfig
}

func (s *DefaultSender) SendLoginCode(_ context.Context, to Email, code LoginCode) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", string(s.config.FromEmail))
	mail.SetHeader("To", string(to))
	mail.SetHeader("Subject", fmt.Sprintf("%s: %d", subject, code))
	mail.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")

	dialer := gomail.NewDialer(
		s.config.Host,
		int(s.config.Port),
		s.config.User,
		s.config.Password,
	)
	return dialer.DialAndSend(mail)
}
