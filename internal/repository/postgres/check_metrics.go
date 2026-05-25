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
