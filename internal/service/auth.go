package service

import (
	"context"
	"errors"
	"strings"

	"github.com/sagemyrage/code-quality-expert-system/internal/domain"
	"github.com/sagemyrage/code-quality-expert-system/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(context.Context, string, string) (*domain.User, error)
	FindByEmail(context.Context, string) (*domain.User, error)
}

type SessionRepository interface {
	Create(context.Context, int64) (string, error)
	GetUserID(context.Context, string) (int64, error)
	Delete(context.Context, string) error
}

type AuthService struct {
	userRepo    UserRepository
	sessionRepo SessionRepository
}

func NewAuthService(userRepo UserRepository, sessionRepo SessionRepository) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

type LoginResult struct {
	UserID    int64
	Email     string
	SessionID string
}

func (s *AuthService) Register(
	ctx context.Context,
	email string,
	password string,
	passwordConfirmation string,
) (*domain.User, error) {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	if email == "" {
		return nil, &ValidationError{Message: "email is required"}
	}

	if password == "" {
		return nil, &ValidationError{Message: "password is required"}
	}
	if len(password) < 8 {
		return nil, &ValidationError{Message: "password must be at least 8 characters"}
	}
	if password != passwordConfirmation {
		return nil, &ValidationError{Message: "passwords do not match"}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	passwordHash := string(hash)

	user, err := s.userRepo.Create(ctx, email, passwordHash)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return nil, &ValidationError{Message: "email already exists"}
		}
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (*LoginResult, error) {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	if email == "" {
		return nil, &ValidationError{Message: "email is required"}
	}

	if password == "" {
		return nil, &ValidationError{Message: "password is required"}
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, &ValidationError{Message: "invalid email or password"}
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, &ValidationError{Message: "invalid email or password"}
	}

	sessionID, err := s.sessionRepo.Create(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		UserID:    user.ID,
		Email:     user.Email,
		SessionID: sessionID,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return nil
	}

	return s.sessionRepo.Delete(ctx, sessionID)
}
