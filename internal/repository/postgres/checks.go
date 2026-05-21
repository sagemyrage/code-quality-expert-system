package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sagemyrage/code-quality-expert-system/internal/domain"
	"github.com/sagemyrage/code-quality-expert-system/internal/repository"
)

type CheckRepository struct {
	pool *pgxpool.Pool
}

func NewCheckRepository(pool *pgxpool.Pool) *CheckRepository {
	return &CheckRepository{
		pool: pool,
	}
}

func (r *CheckRepository) Create(ctx context.Context, userID int64, sourceCode string) (*domain.Check, error) {
	query := `
		INSERT INTO checks (user_id, source_code)
		VALUES ($1, $2)
		RETURNING id, user_id, source_code, created_at, updated_at
	`

	var check domain.Check
	err := r.pool.QueryRow(ctx, query, userID, sourceCode).Scan(
		&check.ID,
		&check.UserID,
		&check.SourceCode,
		&check.CreatedAt,
		&check.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &check, nil
}

func (r *CheckRepository) CreateWithMetrics(ctx context.Context, userID int64, sourceCode string, metrics domain.CheckMetrics) (*domain.Check, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	createCheck := `
		INSERT INTO checks (user_id, source_code)
		VALUES ($1, $2)
		RETURNING id, user_id, source_code, created_at, updated_at
	`

	var check domain.Check
	err = tx.QueryRow(ctx, createCheck, userID, sourceCode).Scan(
		&check.ID,
		&check.UserID,
		&check.SourceCode,
		&check.CreatedAt,
		&check.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	metrics.CheckID = check.ID

	createCheckMetrics := `
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
	err = tx.QueryRow(
		ctx,
		createCheckMetrics,
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

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &check, nil
}

func (r *CheckRepository) ListByUserID(ctx context.Context, userID int64, limit int) ([]domain.Check, error) {
	query := `
		SELECT id, user_id, source_code, created_at, updated_at
		FROM checks
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checks []domain.Check
	for rows.Next() {
		check := domain.Check{}
		if err := rows.Scan(
			&check.ID,
			&check.UserID,
			&check.SourceCode,
			&check.CreatedAt,
			&check.UpdatedAt,
		); err != nil {
			return nil, err
		}

		checks = append(checks, check)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return checks, nil
}

func (r *CheckRepository) FindByIDAndUserID(ctx context.Context, checkID int64, userID int64) (*domain.Check, error) {
	query := `
		SELECT id, user_id, source_code, created_at, updated_at
		FROM checks
		WHERE id = $1 AND user_id = $2
	`

	var check domain.Check
	err := r.pool.QueryRow(ctx, query, checkID, userID).Scan(
		&check.ID,
		&check.UserID,
		&check.SourceCode,
		&check.CreatedAt,
		&check.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrCheckNotFound
		}

		return nil, err
	}

	return &check, nil
}
