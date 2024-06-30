package auth

import (
	"bytes"
	"context"
	"fmt"
	"gopkg.in/gomail.v2"
	frontendEmail "quizzly/web/frontend/templ/email"
)

const (
	subject = "ваш код для входа"
)

type DefaultSender struct {
	config *SenderConfig
}

func (s *DefaultSender) SendLoginCode(ctx context.Context, to Email, code LoginCode) error {
	bodyBuffer := new(bytes.Buffer)
	err := frontendEmail.Code(int(code)).Render(ctx, bodyBuffer)
	if err != nil {
		return err
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", string(s.config.FromEmail))
	mail.SetHeader("To", string(to))
	mail.SetHeader("Subject", fmt.Sprintf("%d — %s", code, subject))
	mail.SetBody("text/html", bodyBuffer.String())

	dialer := gomail.NewDialer(
		s.config.Host,
		int(s.config.Port),
		s.config.User,
		s.config.Password,
	)
	return dialer.DialAndSend(mail)
}
