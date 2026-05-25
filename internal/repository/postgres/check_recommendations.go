package postgres

import (
	"context"

	"github.com/sagemyrage/code-quality-expert-system/internal/domain"
)

func (r *CheckRepository) FindRecommendationsByCheckID(ctx context.Context, checkID int64) ([]domain.CheckRecommendation, error) {
	query := `
	SELECT id, check_id, text, position
	FROM check_recommendations
	WHERE check_id = $1
	ORDER BY position
	`

	rows, err := r.pool.Query(ctx, query, checkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checkRecommendations []domain.CheckRecommendation
	for rows.Next() {
		checkRecommendation := domain.CheckRecommendation{}
		if err := rows.Scan(
			&checkRecommendation.ID,
			&checkRecommendation.CheckID,
			&checkRecommendation.Text,
			&checkRecommendation.Position,
		); err != nil {
			return nil, err
		}

		checkRecommendations = append(checkRecommendations, checkRecommendation)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return checkRecommendations, nil
}
