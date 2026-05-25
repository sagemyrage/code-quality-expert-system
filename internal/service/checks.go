package service

import (
	"context"
	"strings"

	"github.com/sagemyrage/code-quality-expert-system/internal/analyzer"
	"github.com/sagemyrage/code-quality-expert-system/internal/domain"
	"github.com/sagemyrage/code-quality-expert-system/internal/fuzzy"
)

type CheckRepository interface {
	ListByUserID(ctx context.Context, userID int64) ([]domain.Check, error)
	FindByIDAndUserID(ctx context.Context, checkID int64, userID int64) (*domain.Check, error)
	Create(ctx context.Context, userID int64, sourceCode string, score float64, level string, recommendations []string, metrics domain.CheckMetrics) (*domain.Check, error)
	FindMetricsByCheckID(ctx context.Context, checkID int64) (*domain.CheckMetrics, error)
	FindRecommendationsByCheckID(ctx context.Context, checkID int64) ([]domain.CheckRecommendation, error)
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

	fuzzyInput := fuzzy.Input{
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
	evaluation, err := fuzzy.Evaluate(fuzzyInput)
	if err != nil {
		return nil, err
	}

	check, err := s.checkRepo.Create(ctx, userID, sourceCode, evaluation.Score, evaluation.Level, evaluation.Recommendations, checkMetrics)
	if err != nil {
		return nil, err
	}

	return check, nil
}

func (s *CheckService) ListByUserID(ctx context.Context, userID int64) ([]domain.Check, error) {
	return s.checkRepo.ListByUserID(ctx, userID)
}

func (s *CheckService) GetByID(ctx context.Context, checkID int64, userID int64) (*domain.Check, error) {
	return s.checkRepo.FindByIDAndUserID(ctx, checkID, userID)
}

func (s *CheckService) GetDetailsByID(ctx context.Context, checkID int64, userID int64) (*domain.CheckDetails, error) {
	check, err := s.checkRepo.FindByIDAndUserID(ctx, checkID, userID)
	if err != nil {
		return nil, err
	}

	metrics, err := s.checkRepo.FindMetricsByCheckID(ctx, checkID)
	if err != nil {
		return nil, err
	}

	recommendations, err := s.checkRepo.FindRecommendationsByCheckID(ctx, checkID)
	if err != nil {
		return nil, err
	}

	return &domain.CheckDetails{Check: *check, Metrics: *metrics, Recommendations: recommendations}, nil
}
