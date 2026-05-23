package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/sagemyrage/code-quality-expert-system/internal/domain"
	"github.com/sagemyrage/code-quality-expert-system/internal/repository"
)

func (r *CheckRepository) FindMetricsByCheckID(ctx context.Context, checkID int64) (*domain.CheckMetrics, error) {
	query := `
	SELECT
		check_id,
		line_count,
		comment_line_count,
		comment_ratio,
		function_count,
		average_function_length,
		max_function_length,
		conditional_count,
		loop_count,
		max_nesting_depth,
		global_variable_count,
		long_line_count,
		created_at
	FROM check_metrics
	WHERE check_id = $1
	`

	var checkMetrics domain.CheckMetrics
	err := r.pool.QueryRow(ctx, query, checkID).Scan(
		&checkMetrics.CheckID,
		&checkMetrics.LineCount,
		&checkMetrics.CommentLineCount,
		&checkMetrics.CommentRatio,
		&checkMetrics.FunctionCount,
		&checkMetrics.AverageFunctionLength,
		&checkMetrics.MaxFunctionLength,
		&checkMetrics.ConditionalCount,
		&checkMetrics.LoopCount,
		&checkMetrics.MaxNestingDepth,
		&checkMetrics.GlobalVariableCount,
		&checkMetrics.LongLineCount,
		&checkMetrics.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrCheckMetricsNotFound
		}

		return nil, err
	}

	return &checkMetrics, nil
}

func (r *CheckRepository) CreateMetrics(ctx context.Context, metrics domain.CheckMetrics) (*domain.CheckMetrics, error) {
	query := `
	INSERT INTO check_metrics (
		check_id,
		line_count,
		comment_line_count,
		comment_ratio,
		function_count,
		average_function_length,
		max_function_length,
		conditional_count,
		loop_count,
		max_nesting_depth,
		global_variable_count,
		long_line_count
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	RETURNING 
		check_id,
		line_count,
		comment_line_count,
		comment_ratio,
		function_count,
		average_function_length,
		max_function_length,
		conditional_count,
		loop_count,
		max_nesting_depth,
		global_variable_count,
		long_line_count,
		created_at
	`

	var createdMetrics domain.CheckMetrics
	err := r.pool.QueryRow(
		ctx,
		query,
		metrics.CheckID,
		metrics.LineCount,
		metrics.CommentLineCount,
		metrics.CommentRatio,
		metrics.FunctionCount,
		metrics.AverageFunctionLength,
		metrics.MaxFunctionLength,
		metrics.ConditionalCount,
		metrics.LoopCount,
		metrics.MaxNestingDepth,
		metrics.GlobalVariableCount,
		metrics.LongLineCount,
	).Scan(
		&createdMetrics.CheckID,
		&createdMetrics.LineCount,
		&createdMetrics.CommentLineCount,
		&createdMetrics.CommentRatio,
		&createdMetrics.FunctionCount,
		&createdMetrics.AverageFunctionLength,
		&createdMetrics.MaxFunctionLength,
		&createdMetrics.ConditionalCount,
		&createdMetrics.LoopCount,
		&createdMetrics.MaxNestingDepth,
		&createdMetrics.GlobalVariableCount,
		&createdMetrics.LongLineCount,
		&createdMetrics.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &createdMetrics, nil
}
