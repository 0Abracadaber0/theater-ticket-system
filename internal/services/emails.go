package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
	"theater-ticket-system/internal/config"
)

type EmailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{cfg: cfg}
}

// GenerateCode создает 6-значный код
func (s *EmailService) GenerateCode() (string, error) {
	max := big.NewInt(999999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()%1000000), nil
}

// SendVerificationCode отправляет код подтверждения на email
func (s *EmailService) SendVerificationCode(email, code string) error {
	from := s.cfg.Email.From
	password := s.cfg.Email.Password
	smtpHost := s.cfg.Email.SMTPHost
	smtpPort := s.cfg.Email.SMTPPort

	// Формируем сообщение
	subject := "Код подтверждения - Театральная касса"
	body := fmt.Sprintf(`
Здравствуйте!

Ваш код подтверждения: %s

Код действителен в течение 10 минут.

Если вы не запрашивали этот код, просто проигнорируйте это письмо.

--
Театральная касса
`, code)

	message := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", from, email, subject, body))

	// Аутентификация
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Отправка письма
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	err := smtp.SendMail(addr, auth, from, []string{email}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
