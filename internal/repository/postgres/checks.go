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
