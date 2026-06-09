package repository

import (
	"context"

	"github.com/Touy2004/palm-back-end/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type PalmRepository struct {
	db *pgxpool.Pool
}

func NewPalmRepository(db *pgxpool.Pool) *PalmRepository {
	return &PalmRepository{db: db}
}

func (r *PalmRepository) Create(template *model.PalmTemplate) error {
	if template.ID == uuid.Nil {
		template.ID = uuid.New()
	}
	if template.CreatedAt.IsZero() {
		now := time.Now()
		template.CreatedAt = now
		template.UpdatedAt = now
	}

	query := `
		INSERT INTO palm_templates (id, user_id, hand_side, template_encrypted, template_nonce, embedding_dim, model_version, threshold, status, registered_device_id, created_at, updated_at, revoked_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at`
	
	return r.db.QueryRow(context.Background(), query,
		template.ID, template.UserID, template.HandSide, template.TemplateEncrypted,
		template.TemplateNonce, template.EmbeddingDim, template.ModelVersion,
		template.Threshold, template.Status, template.RegisteredDeviceID,
		template.CreatedAt, template.UpdatedAt, template.RevokedAt,
	).Scan(&template.ID, &template.CreatedAt, &template.UpdatedAt)
}

func (r *PalmRepository) FindByUserID(userID string) ([]model.PalmTemplate, error) {
	var templates []model.PalmTemplate
	query := `
		SELECT p.id, p.user_id, p.hand_side, p.template_encrypted, p.template_nonce, p.embedding_dim, p.model_version, p.threshold, p.status, p.registered_device_id, p.created_at, p.updated_at, p.revoked_at,
		       d.name, d.device_code
		FROM palm_templates p
		LEFT JOIN devices d ON p.registered_device_id = d.id
		WHERE p.user_id = $1 AND p.status = 'active'
	`
	
	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t model.PalmTemplate
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.HandSide, &t.TemplateEncrypted, &t.TemplateNonce,
			&t.EmbeddingDim, &t.ModelVersion, &t.Threshold, &t.Status, &t.RegisteredDeviceID,
			&t.CreatedAt, &t.UpdatedAt, &t.RevokedAt,
			&t.RegisteredDeviceName, &t.RegisteredDeviceCode,
		); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, rows.Err()
}

func (r *PalmRepository) FindAllActive() ([]model.PalmTemplate, error) {
	var templates []model.PalmTemplate
	query := `SELECT id, user_id, hand_side, template_encrypted, template_nonce, embedding_dim, model_version, threshold, status, registered_device_id, created_at, updated_at, revoked_at FROM palm_templates WHERE status = 'active'`
	
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t model.PalmTemplate
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.HandSide, &t.TemplateEncrypted, &t.TemplateNonce,
			&t.EmbeddingDim, &t.ModelVersion, &t.Threshold, &t.Status, &t.RegisteredDeviceID,
			&t.CreatedAt, &t.UpdatedAt, &t.RevokedAt,
		); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, rows.Err()
}

func (r *PalmRepository) Delete(id, userID string) error {
	query := `UPDATE palm_templates SET status = 'revoked', revoked_at = NOW() WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(context.Background(), query, id, userID)
	return err
}