package domain

import "context"

type Mailer interface {
	SendPasswordResetEmail(ctx context.Context, toEmail, resetToken string) error
	SendActivationEmail(ctx context.Context, toEmail, activationToken string) error
}
