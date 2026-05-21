package service

import (
	"context"
	"strings"

	"github.com/sagemyrage/code-quality-expert-system/internal/analyzer"
	"github.com/sagemyrage/code-quality-expert-system/internal/domain"
)

const checkHistoryLimit = 10

type CheckRepository interface {
	ListByUserID(context.Context, int64, int) ([]domain.Check, error)
	FindByIDAndUserID(context.Context, int64, int64) (*domain.Check, error)
	CreateWithMetrics(context.Context, int64, string, domain.CheckMetrics) (*domain.Check, error)
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

	analysisMetrics, err := analyzer.Analyze(sourceCode)
	if err != nil {
		return nil, &ValidationError{Message: "source code must be valid Go code"}
	}

	checkMetrics := domain.CheckMetrics{
		LineCount:             analysisMetrics.LineCount,
		CommentLineCount:      analysisMetrics.CommentLineCount,
		CommentRatio:          analysisMetrics.CommentRatio,
		FunctionCount:         analysisMetrics.FunctionCount,
		AverageFunctionLength: analysisMetrics.AverageFunctionLength,
		MaxFunctionLength:     analysisMetrics.MaxFunctionLength,
		ConditionalCount:      analysisMetrics.ConditionalCount,
		LoopCount:             analysisMetrics.LoopCount,
		MaxNestingDepth:       analysisMetrics.MaxNestingDepth,
		GlobalVariableCount:   analysisMetrics.GlobalVariableCount,
		LongLineCount:         analysisMetrics.LongLineCount,
	}
	check, err := s.checkRepo.CreateWithMetrics(ctx, userID, sourceCode, checkMetrics)
	if err != nil {
		return nil, err
	}

	return check, nil
}

func (s *CheckService) ListByUserID(ctx context.Context, userID int64) ([]domain.Check, error) {
	return s.checkRepo.ListByUserID(ctx, userID, checkHistoryLimit)
}

func (s *CheckService) GetByID(ctx context.Context, checkID int64, userID int64) (*domain.Check, error) {
	return s.checkRepo.FindByIDAndUserID(ctx, checkID, userID)
}
