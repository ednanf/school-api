package mailer

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/ednanf/school-api/internal/domain"
)

type MailpitService struct {
	host     string
	port     int
	fromAddr string
}

func NewMailpitService(host string, port int, fromAddr string) domain.Mailer {
	return &MailpitService{
		host:     host,
		port:     port,
		fromAddr: fromAddr,
	}
}

func (m *MailpitService) SendPasswordResetEmail(ctx context.Context, toEmail, resetToken string) error {
	addr := fmt.Sprintf("%s:%d", m.host, m.port)

	subject := "Subject: Reset Your Password\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	// Link points to frontend or local test URL
	resetLink := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", resetToken)

	body := fmt.Sprintf(`
		<html>
			<body>
				<h2>Password Reset Request</h2>
				<p>You requested a password reset for your account.</p>
				<p>Click the link below to set a new password:</p>
				<p><a href="%s">%s</a></p>
				<br>
				<p>If you did not request this, please ignore this email.</p>
			</body>
		</html>
	`, resetLink, resetLink)

	msg := []byte(subject + mime + body)

	// Mailpit runs locally without authentication by default
	err := smtp.SendMail(addr, nil, m.fromAddr, []string{toEmail}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email via Mailpit: %w", err)
	}

	return nil
}

// SendActivationEmail sends a raw HTML activation email via local Mailpit SMTP
func (m *MailpitService) SendActivationEmail(ctx context.Context, toEmail, activationToken string) error {
	activationLink := fmt.Sprintf("http://localhost:3000/activate?token=%s", activationToken)

	// Format MIME headers and HTML body
	header := fmt.Sprintf("From: %s\r\n", m.fromAddr) +
		fmt.Sprintf("To: %s\r\n", toEmail) +
		"Subject: Activate Your Account\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n\r\n"

	body := fmt.Sprintf(`
		<html>
			<body>
				<h2>Welcome to School API!</h2>
				<p>An account has been created for you.</p>
				<p>Please click the link below to activate your account:</p>
				<p><a href="%s">%s</a></p>
			</body>
		</html>
	`, activationLink, activationLink)

	msg := []byte(header + body)
	addr := fmt.Sprintf("%s:%d", m.host, m.port)

	// Mailpit doesn't require SMTP auth in local development
	err := smtp.SendMail(addr, nil, m.fromAddr, []string{toEmail}, msg)
	if err != nil {
		return fmt.Errorf("failed to send activation email via Mailpit: %w", err)
	}

	return nil
}
