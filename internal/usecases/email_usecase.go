package usecases

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/errors"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

// EmailUseCase defines the interface for email operations
type EmailUseCase interface {
	SendPasswordResetEmail(to, resetLink string) error
}

// emailUseCase implements EmailUseCase interface
type emailUseCase struct {
	dialer *gomail.Dialer
	from   string
	logger *zap.Logger
}

// NewEmailUseCase creates a new instance of EmailUseCase
func NewEmailUseCase(logger *zap.Logger) (EmailUseCase, error) {
	host := os.Getenv("MAIL_HOST")
	portStr := os.Getenv("MAIL_PORT")
	username := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	from := os.Getenv("MAIL_FROM")

	logger.Debug("Loading email configuration from environment variables",
		zap.String("host", host),
		zap.String("port", portStr),
		zap.String("username", username),
		zap.String("from", from),
	)

	if host == "" {
		logger.Error("MAIL_HOST environment variable is missing")
		return nil, errors.WrapError(errors.ErrEmailConfigMissing, "MAIL_HOST environment variable is required")
	}
	if portStr == "" {
		logger.Error("MAIL_PORT environment variable is missing")
		return nil, errors.WrapError(errors.ErrEmailConfigMissing, "MAIL_PORT environment variable is required")
	}
	if username == "" {
		logger.Error("MAIL_USERNAME environment variable is missing")
		return nil, errors.WrapError(errors.ErrEmailConfigMissing, "MAIL_USERNAME environment variable is required")
	}
	if password == "" {
		logger.Error("MAIL_PASSWORD environment variable is missing")
		return nil, errors.WrapError(errors.ErrEmailConfigMissing, "MAIL_PASSWORD environment variable is required")
	}
	if from == "" {
		logger.Error("MAIL_FROM environment variable is missing")
		return nil, errors.WrapError(errors.ErrEmailConfigMissing, "MAIL_FROM environment variable is required")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		logger.Error("Failed to parse MAIL_PORT",
			zap.String("port", portStr),
			zap.Error(err),
		)
		return nil, errors.WrapError(errors.ErrInvalidPort, fmt.Sprintf("invalid MAIL_PORT: %s", portStr))
	}

	dialer := gomail.NewDialer(host, port, username, password)

	logger.Info("Email use case initialized successfully",
		zap.String("host", host),
		zap.Int("port", port),
		zap.String("username", username),
		zap.String("from", from),
	)

	return &emailUseCase{
		dialer: dialer,
		from:   from,
		logger: logger,
	}, nil
}

// SendPasswordResetEmail sends a password reset email with magic link
func (uc *emailUseCase) SendPasswordResetEmail(to, resetLink string) error {
	uc.logger.Info("Sending password reset email",
		zap.String("to", to),
		zap.String("reset_link", resetLink),
	)

	subject := "Recovery Password - DafonCV"

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Recovery Password</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					line-height: 1.6;
					color: #333;
					max-width: 600px;
					margin: 0 auto;
					padding: 20px;
				}
				.header {
					background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
					color: white;
					padding: 30px;
					text-align: center;
					border-radius: 10px 10px 0 0;
				}
				.content {
					background: #f9f9f9;
					padding: 30px;
					border-radius: 0 0 10px 10px;
				}
				.button {
					display: inline-block;
					background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
					color: white;
					padding: 15px 30px;
					text-decoration: none;
					border-radius: 5px;
					margin: 20px 0;
					font-weight: bold;
				}
				.footer {
					margin-top: 30px;
					padding-top: 20px;
					border-top: 1px solid #ddd;
					font-size: 12px;
					color: #666;
				}
				.warning {
					background: #fff3cd;
					border: 1px solid #ffeaa7;
					padding: 15px;
					border-radius: 5px;
					margin: 20px 0;
				}
			</style>
		</head>
		<body>
			<div class="header">
				<h1>🔐 Recovery Password</h1>
				<p>DafonCV - Curriculum System</p>
			</div>
			
			<div class="content">
				<h2>Hello!</h2>
				<p>We received a request to reset your password on the DafonCV platform.</p>
				
				<p>If you did not request this change, you can safely ignore this email.</p>
				
				<div style="text-align: center;">
					<a href="%s" class="button">🔑 Reset My Password</a>
				</div>
				
				<div class="warning">
					<strong>⚠️ Important:</strong>
					<ul>
						<li>This link is valid for 1 hour</li>
						<li>Use it only once</li>
						<li>Do not share this link with anyone</li>
					</ul>
				</div>
				
				<p>If the button doesn't work, copy and paste the link below into your browser:</p>
				<p style="word-break: break-all; background: #f8f9fa; padding: 10px; border-radius: 5px; font-size: 12px;">
					%s
				</p>
			</div>
			
			<div class="footer">
				<p>This is an automated email, do not reply to this message.</p>
				<p>© 2024 DafonCV. All rights reserved.</p>
			</div>
		</body>
		</html>
	`, resetLink, resetLink)

	textBody := fmt.Sprintf(`
Reset Password - DafonCV

Hello!

We received a request to reset your password on the DafonCV platform.

If you did not request this change, you can safely ignore this email.

To reset your password, access the link:
%s

⚠️ Important:
- This link is valid for 1 hour
- Use it only once
- Do not share this link with anyone

This is an automated email, do not reply to this message.

© 2025 DafonCV. All rights reserved.
	`, resetLink)

	m := gomail.NewMessage()
	m.SetHeader("From", uc.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", textBody)
	m.AddAlternative("text/html", htmlBody)

	if err := uc.dialer.DialAndSend(m); err != nil {
		uc.logger.Error("Failed to send password reset email",
			zap.String("to", to),
			zap.Error(err),
		)
		return errors.WrapError(errors.ErrEmailSendFailed, fmt.Sprintf("failed to send email to %s", to))
	}

	uc.logger.Info("Password reset email sent successfully",
		zap.String("to", to),
	)

	return nil
}

// GenerateSecureToken generates a secure random token for password reset
func GenerateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
