package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sagemyrage/code-quality-expert-system/internal/domain"
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
