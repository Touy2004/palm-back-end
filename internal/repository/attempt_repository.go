package repository

import (
	"context"

	"github.com/Touy2004/palm-back-end/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type AttemptRepository struct {
	db *pgxpool.Pool
}

func NewAttemptRepository(db *pgxpool.Pool) *AttemptRepository {
	return &AttemptRepository{db: db}
}

func (r *AttemptRepository) Create(attempt *model.PalmAuthAttempt) error {
	if attempt.ID == uuid.Nil {
		attempt.ID = uuid.New()
	}
	if attempt.CreatedAt.IsZero() {
		attempt.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO palm_auth_attempts (id, user_id, device_id, template_id, action, score, threshold, liveness_passed, quality_score, thermal_min, thermal_max, thermal_avg, result, failure_reason, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, created_at`
	
	return r.db.QueryRow(context.Background(), query,
		attempt.ID, attempt.UserID, attempt.DeviceID, attempt.TemplateID, attempt.Action,
		attempt.Score, attempt.Threshold, attempt.LivenessPassed, attempt.QualityScore,
		attempt.ThermalMin, attempt.ThermalMax, attempt.ThermalAvg, attempt.Result, attempt.FailureReason, attempt.CreatedAt,
	).Scan(&attempt.ID, &attempt.CreatedAt)
}
