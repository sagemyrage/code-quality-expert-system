package service

import (
	"context"
	"strings"

	"github.com/sagemyrage/code-quality-expert-system/internal/domain"
)

const checkHistoryLimit = 10

type CheckRepository interface {
	Create(context.Context, int64, string) (*domain.Check, error)
	ListByUserID(context.Context, int64, int) ([]domain.Check, error)
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

func (s *CheckService) ListByUserID(ctx context.Context, userID int64) ([]domain.Check, error) {
	return s.checkRepo.ListByUserID(ctx, userID, checkHistoryLimit)
}
