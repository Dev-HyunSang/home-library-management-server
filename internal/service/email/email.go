package email

import (
	"fmt"
	"net/smtp"

	"github.com/dev-hyunsang/home-library-backend/internal/config"
	"github.com/dev-hyunsang/home-library-backend/logger"
)

// Service handles email sending operations
type Service struct {
	smtpHost     string
	smtpPort     string
	fromAddress  string
	fromPassword string
}

// Config holds email service configuration
type Config struct {
	SMTPHost     string
	SMTPPort     string
	FromAddress  string
	FromPassword string
}

// NewService creates a new email service with explicit configuration
func NewService(cfg Config) *Service {
	return &Service{
		smtpHost:     cfg.SMTPHost,
		smtpPort:     cfg.SMTPPort,
		fromAddress:  cfg.FromAddress,
		fromPassword: cfg.FromPassword,
	}
}

// NewServiceFromEnv creates an email service using environment variables
func NewServiceFromEnv() *Service {
	return &Service{
		smtpHost:     config.GetEnv("GOOGLE_MAIL_SMTP"),
		smtpPort:     config.DefaultSMTPPort,
		fromAddress:  config.GetEnv("GOOGLE_MAIL_ADDRESS"),
		fromPassword: config.GetEnv("GOOGLE_MAIL_PASSWORD"),
	}
}

func (s *Service) smtpAddr() string {
	return fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)
}

func (s *Service) auth() smtp.Auth {
	return smtp.PlainAuth("", s.fromAddress, s.fromPassword, s.smtpHost)
}

// SendVerificationCode sends an email verification code
func (s *Service) SendVerificationCode(toEmail, code string) error {
	subject := "Subject: 나만의 서재 메일 인증번호입니다.\r\n"
	body := fmt.Sprintf("인증번호는 %s입니다.\r\n", code)

	return s.send(toEmail, subject, body)
}

// SendPasswordReset sends a temporary password
func (s *Service) SendPasswordReset(toEmail, tempPassword string) error {
	subject := "Subject: 비밀번호 초기화\r\n"
	body := fmt.Sprintf("비밀번호 초기화 메일 테스트입니다. %s\r\n", tempPassword)

	return s.send(toEmail, subject, body)
}

func (s *Service) send(toEmail, subject, body string) error {
	msg := []byte(subject + "\r\n" + body)

	err := smtp.SendMail(s.smtpAddr(), s.auth(), s.fromAddress, []string{toEmail}, msg)
	if err != nil {
		logger.Sugar().Errorf("이메일 발송 실패 (to: %s): %v", toEmail, err)
		return fmt.Errorf("이메일 발송 실패: %w", err)
	}

	logger.Sugar().Infof("이메일 발송 성공: %s", toEmail)
	return nil
}
