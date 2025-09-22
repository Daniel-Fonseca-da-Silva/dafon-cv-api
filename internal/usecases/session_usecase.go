package usecases

import (
	"time"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/errors"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SessionUseCase interface {
	CreateSession(userID uuid.UUID, duration time.Duration) (*models.Session, error)
	GetSessionByToken(token string) (*models.Session, error)
	GetActiveSessionsByUserID(userID uuid.UUID) ([]*models.Session, error)
	ValidateSession(token string) (*models.Session, error)
	DeactivateSession(token string) error
	DeactivateUserSessions(userID uuid.UUID) error
	CleanupExpiredSessions() error
}

type sessionUseCase struct {
	sessionRepo repositories.SessionRepository
	emailUC     EmailUseCase
	logger      *zap.Logger
}

func NewSessionUseCase(sessionRepo repositories.SessionRepository, emailUC EmailUseCase, logger *zap.Logger) SessionUseCase {
	return &sessionUseCase{
		sessionRepo: sessionRepo,
		emailUC:     emailUC,
		logger:      logger,
	}
}

func (uc *sessionUseCase) CreateSession(userID uuid.UUID, duration time.Duration) (*models.Session, error) {
	uc.logger.Info("Creating new session",
		zap.String("user_id", userID.String()),
		zap.Duration("duration", duration),
	)

	// Generate secure token
	token, err := GenerateSecureToken()
	if err != nil {
		uc.logger.Error("Failed to generate secure token",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		return nil, errors.WrapError(errors.ErrTokenGenerationFailed, "failed to generate secure token")
	}

	// Create session with configurable expiration
	session := &models.Session{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(duration),
		IsActive:  true,
	}

	if err := uc.sessionRepo.Create(session); err != nil {
		uc.logger.Error("Failed to create session",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		return nil, errors.WrapError(errors.ErrSessionCreationFailed, "failed to create session")
	}

	uc.logger.Info("Session created successfully",
		zap.String("user_id", userID.String()),
		zap.String("session_id", session.ID.String()),
		zap.Time("expires_at", session.ExpiresAt),
	)

	return session, nil
}

func (uc *sessionUseCase) GetSessionByToken(token string) (*models.Session, error) {
	session, err := uc.sessionRepo.GetByToken(token)
	if err != nil {
		uc.logger.Error("Failed to get session by token",
			zap.String("token", token),
			zap.Error(err),
		)
		return nil, errors.WrapError(errors.ErrSessionNotFound, "session not found")
	}
	return session, nil
}

func (uc *sessionUseCase) GetActiveSessionsByUserID(userID uuid.UUID) ([]*models.Session, error) {
	sessions, err := uc.sessionRepo.GetActiveByUserID(userID)
	if err != nil {
		uc.logger.Debug("No active sessions found for user",
			zap.String("user_id", userID.String()),
		)
		return nil, errors.WrapError(errors.ErrSessionNotFound, "no active sessions found")
	}
	return sessions, nil
}

func (uc *sessionUseCase) ValidateSession(token string) (*models.Session, error) {
	session, err := uc.GetSessionByToken(token)
	if err != nil {
		return nil, err
	}

	if !session.IsValid() {
		uc.logger.Warn("Session is invalid or expired",
			zap.String("session_id", session.ID.String()),
			zap.String("user_id", session.UserID.String()),
			zap.Bool("is_active", session.IsActive),
			zap.Time("expires_at", session.ExpiresAt),
		)
		return nil, errors.WrapError(errors.ErrSessionExpired, "session is expired or inactive")
	}

	return session, nil
}

func (uc *sessionUseCase) DeactivateSession(token string) error {
	uc.logger.Info("Deactivating session",
		zap.String("token", token),
	)

	if err := uc.sessionRepo.DeactivateByToken(token); err != nil {
		uc.logger.Error("Failed to deactivate session",
			zap.String("token", token),
			zap.Error(err),
		)
		return errors.WrapError(errors.ErrSessionDeactivationFailed, "failed to deactivate session")
	}

	uc.logger.Info("Session deactivated successfully",
		zap.String("token", token),
	)

	return nil
}

func (uc *sessionUseCase) DeactivateUserSessions(userID uuid.UUID) error {
	uc.logger.Info("Deactivating all sessions for user",
		zap.String("user_id", userID.String()),
	)

	if err := uc.sessionRepo.DeactivateByUserID(userID); err != nil {
		uc.logger.Error("Failed to deactivate user sessions",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		return errors.WrapError(errors.ErrSessionDeactivationFailed, "failed to deactivate user sessions")
	}

	uc.logger.Info("User sessions deactivated successfully",
		zap.String("user_id", userID.String()),
	)

	return nil
}

func (uc *sessionUseCase) CleanupExpiredSessions() error {
	uc.logger.Info("Cleaning up expired sessions")

	if err := uc.sessionRepo.DeleteExpired(); err != nil {
		uc.logger.Error("Failed to cleanup expired sessions",
			zap.Error(err),
		)
		return errors.WrapError(errors.ErrSessionCleanupFailed, "failed to cleanup expired sessions")
	}

	uc.logger.Info("Expired sessions cleaned up successfully")
	return nil
}
