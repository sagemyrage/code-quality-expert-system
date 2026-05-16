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
}

type AuthService struct {
	userRepo UserRepository
}

func NewAuthService(userRepo UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
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
