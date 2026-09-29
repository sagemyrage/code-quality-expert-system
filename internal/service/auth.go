package service

import (
	"context"
	"errors"
	"regexp"
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

var emailRegex = regexp.MustCompile(
	`^[a-z0-9](?:[a-z0-9_-]*[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9_-]*[a-z0-9])?)*@` +
		`[a-z0-9](?:[a-z0-9-]*[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)+$`,
)

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

	if !isValidEmail(email) {
		return nil, &ValidationError{Message: "invalid email"}
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

	if !isValidEmail(email) {
		return nil, &ValidationError{Message: "invalid email"}
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

func (s *AuthService) AuthenticateSession(ctx context.Context, sessionID string) (int64, error) {
	if sessionID == "" {
		return 0, ErrUnauthenticated
	}

	userID, err := s.sessionRepo.GetUserID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, repository.ErrSessionNotFound) {
			return 0, ErrUnauthenticated
		}
		return 0, err
	}

	return userID, nil
}

func isValidEmail(email string) bool {
	if len(email) > 254 {
		return false
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	localPart := parts[0]
	domain := parts[1]

	if len(localPart) > 64 || len(localPart) == 0 {
		return false
	}

	for _, label := range strings.Split(domain, ".") {
		if len(label) == 0 || len(label) > 63 {
			return false
		}
	}

	return emailRegex.MatchString(email)
}
