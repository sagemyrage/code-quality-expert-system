package service

import (
	"context"
	"strings"

	"github.com/sagemyrage/code-quality-expert-system/internal/domain"
)

type CheckRepository interface {
	Create(context.Context, int64, string) (*domain.Check, error)
}

type CheckService struct {
	checkRepo CheckRepository
}

func NewCheckService(checkRepo CheckRepository) *CheckService {
	return &CheckService{
		checkRepo: checkRepo,
	}
}

func (s *CheckService) Create(ctx context.Context, userID int64, sourceCode string) (*domain.Check, error) {
	sourceCode = strings.TrimSpace(sourceCode)
	if sourceCode == "" {
		return nil, &ValidationError{Message: "source code is required"}
	}

	return s.checkRepo.Create(ctx, userID, sourceCode)
}
