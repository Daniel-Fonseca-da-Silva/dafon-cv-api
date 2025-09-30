package usecases

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/errors"
	"github.com/resend/resend-go/v2"
	"go.uber.org/zap"
)

// EmailUseCase defines the interface for email operations
type EmailUseCase interface {
	SendSessionTokenEmail(to, name, token string) error
}

// emailUseCase implements EmailUseCase interface
type emailUseCase struct {
	client *resend.Client
	from   string
	logger *zap.Logger
}

// NewEmailUseCase creates a new instance of EmailUseCase
func NewEmailUseCase(logger *zap.Logger) (EmailUseCase, error) {
	apiKey := os.Getenv("RESEND_API_KEY")
	from := os.Getenv("MAIL_FROM")

	logger.Debug("Loading email configuration from environment variables",
		zap.String("from", from),
	)

	if apiKey == "" {
		logger.Error("RESEND_API_KEY environment variable is missing")
		return nil, errors.WrapError(errors.ErrEmailConfigMissing, "RESEND_API_KEY environment variable is required")
	}
	if from == "" {
		logger.Error("MAIL_FROM environment variable is missing")
		return nil, errors.WrapError(errors.ErrEmailConfigMissing, "MAIL_FROM environment variable is required")
	}

	client := resend.NewClient(apiKey)

	logger.Info("Email use case initialized successfully with Resend",
		zap.String("from", from),
	)

	return &emailUseCase{
		client: client,
		from:   from,
		logger: logger,
	}, nil
}

// SendSessionTokenEmail sends a session token to the user's email
func (uc *emailUseCase) SendSessionTokenEmail(to, name, token string) error {
	uc.logger.Info("Sending session token email",
		zap.String("to", to),
		zap.String("name", name),
	)

	// The token is passed directly from the frontend, no need to create a link
	// The frontend will handle the token processing
	loginLink := token

	subject := "Welcome to Dafon CV - Your AI-Powered Resume Builder"
	htmlContent := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Dafon CV - AI Resume Builder</title>
		</head>
		<body style="margin: 0; padding: 0; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; min-height: 100vh;">
			<div style="max-width: 600px; margin: 0 auto; padding: 40px 20px;">
				<!-- Glassmorphism Container -->
				<div style="background: rgba(255, 255, 255, 0.1); backdrop-filter: blur(20px); border-radius: 20px; border: 1px solid rgba(255, 255, 255, 0.2); padding: 40px; box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);">
					
					<!-- Header -->
					<div style="text-align: center; margin-bottom: 40px;">
						<div style="background: rgba(255, 255, 255, 0.2); border-radius: 50%%; width: 80px; height: 80px; margin: 0 auto 20px; display: flex; align-items: center; justify-content: center; border: 2px solid rgba(255, 255, 255, 0.3);">
							<span style="font-size: 32px;">📄</span>
						</div>
						<h1 style="color: white; margin: 0; font-size: 28px; font-weight: 300; text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);">Dafon CV</h1>
						<p style="color: rgba(255, 255, 255, 0.8); margin: 10px 0 0; font-size: 16px; font-weight: 300;">AI-Powered Resume Builder</p>
					</div>

					<!-- Content -->
					<div style="background: rgba(255, 255, 255, 0.95); border-radius: 15px; padding: 30px; margin-bottom: 30px; box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);">
						<h2 style="color: #2c3e50; margin: 0 0 20px; font-size: 24px; font-weight: 400;">Welcome to Your Professional Journey, %s!</h2>
						
						<p style="color: #5a6c7d; line-height: 1.6; margin: 0 0 20px; font-size: 16px;">
							Thank you for choosing <strong style="color: #667eea;">Dafon CV</strong> - the revolutionary platform that transforms your career story into stunning, AI-optimized resumes.
						</p>

						<p style="color: #5a6c7d; line-height: 1.6; margin: 0 0 30px; font-size: 16px;">
							Our intelligent system analyzes your experience and creates personalized, ATS-friendly resumes that stand out to recruiters and hiring managers.
						</p>

						<!-- CTA Button -->
						<div style="text-align: center; margin: 30px 0;">
							<a href="%s" style="background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 16px 32px; text-decoration: none; border-radius: 50px; display: inline-block; font-weight: 600; font-size: 16px; box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4); transition: all 0.3s ease; border: none;">
								🚀 Access Your Dashboard
							</a>
						</div>

						<!-- Features -->
						<div style="background: rgba(102, 126, 234, 0.1); border-radius: 10px; padding: 20px; margin: 30px 0;">
							<h3 style="color: #667eea; margin: 0 0 15px; font-size: 18px; font-weight: 500;">✨ What makes Dafon CV special:</h3>
							<ul style="color: #5a6c7d; margin: 0; padding-left: 20px; line-height: 1.8;">
								<li><strong>AI-Powered Optimization:</strong> Smart content suggestions tailored to your industry</li>
								<li><strong>ATS-Friendly Templates:</strong> Designed to pass Applicant Tracking Systems</li>
								<li><strong>Real-time Analytics:</strong> Track your resume performance and views</li>
								<li><strong>Multiple Formats:</strong> Export to PDF, Word, or share online</li>
							</ul>
						</div>

						<!-- Security Notice -->
						<div style="background: rgba(255, 193, 7, 0.1); border-left: 4px solid #ffc107; padding: 15px; margin: 20px 0; border-radius: 0 8px 8px 0;">
							<p style="color: #856404; margin: 0; font-size: 14px; line-height: 1.5;">
								<strong>🔒 Security Notice:</strong> This secure login link expires in 15 minutes and can only be used once. If you didn't request access to Dafon CV, please ignore this email.
							</p>
						</div>
					</div>

					<!-- Footer -->
					<div style="text-align: center; color: rgba(255, 255, 255, 0.7); font-size: 14px;">
						<p style="margin: 0 0 10px;">Ready to build your dream career?</p>
						<p style="margin: 0; font-size: 12px;">
							© 2025 Dafon CV. Empowering professionals worldwide with AI-driven resume solutions.
						</p>
					</div>
				</div>
			</div>
		</body>
		</html>
	`, name, loginLink)

	params := &resend.SendEmailRequest{
		From:    uc.from,
		To:      []string{to},
		Subject: subject,
		Html:    htmlContent,
	}

	sent, err := uc.client.Emails.Send(params)
	if err != nil {
		uc.logger.Error("Failed to send session token email",
			zap.String("to", to),
			zap.Error(err),
		)
		return errors.WrapError(errors.ErrEmailSendFailed, "failed to send session token email")
	}

	uc.logger.Info("Session token email sent successfully",
		zap.String("to", to),
		zap.String("email_id", sent.Id),
	)

	return nil
}

// GenerateSecureToken generates a secure random token for session
func GenerateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
