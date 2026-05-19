package domain

import "time"

type Check struct {
	ID         int64
	UserID     int64
	SourceCode string
	Score      *float64
	Summary    *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
