package mailer

import (
	"context"
	"fmt"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/resend/resend-go/v4"
)

type ResendService struct {
	client   *resend.Client
	fromAddr string
}

func NewResendService(apiKey, fromAddr string) domain.Mailer {
	return &ResendService{
		client:   resend.NewClient(apiKey),
		fromAddr: fromAddr,
	}
}

func (r *ResendService) SendPasswordResetEmail(ctx context.Context, toEmail, resetToken string) error {
	resetLink := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", resetToken)

	htmlBody := fmt.Sprintf(`
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

	params := &resend.SendEmailRequest{
		From:    r.fromAddr,
		To:      []string{toEmail},
		Subject: "Reset Your Password",
		Html:    htmlBody,
	}

	_, err := r.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email via Resend v4: %w", err)
	}

	return nil
}

func (r *ResendService) SendActivationEmail(ctx context.Context, toEmail, activationToken string) error {
	activationLink := fmt.Sprintf("http://localhost:3000/activate?token=%s", activationToken)

	htmlBody := fmt.Sprintf(`
		<html>
			<body>
				<h2>Welcome to School API!</h2>
				<p>An account has been created for you.</p>
				<p>Please click the link below to activate your account:</p>
				<p><a href="%s">%s</a></p>
			</body>
		</html>
	`, activationLink, activationLink)

	params := &resend.SendEmailRequest{
		From:    r.fromAddr,
		To:      []string{toEmail},
		Subject: "Activate Your Account",
		Html:    htmlBody,
	}

	_, err := r.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send activation email via Resend: %w", err)
	}

	return nil
}
