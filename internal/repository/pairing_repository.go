package repository

import (
	"context"
	"errors"

	"github.com/Touy2004/palm-back-end/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type PairingRepository struct {
	db *pgxpool.Pool
}

func NewPairingRepository(db *pgxpool.Pool) *PairingRepository {
	return &PairingRepository{db: db}
}

func (r *PairingRepository) Create(session *model.DevicePairingSession) error {
	if session.ID == uuid.Nil {
		session.ID = uuid.New()
	}
	if session.CreatedAt.IsZero() {
		session.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO device_pairing_sessions (id, device_id, session_token, user_id, purpose, status, expires_at, scanned_at, approved_at, completed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at`
	
	return r.db.QueryRow(context.Background(), query,
		session.ID, session.DeviceID, session.SessionToken, session.UserID,
		session.Purpose, session.Status, session.ExpiresAt, session.ScannedAt,
		session.ApprovedAt, session.CompletedAt, session.CreatedAt,
	).Scan(&session.ID, &session.CreatedAt)
}

func (r *PairingRepository) FindByID(id string) (*model.DevicePairingSession, error) {
	var session model.DevicePairingSession
	query := `SELECT id, device_id, session_token, user_id, purpose, status, expires_at, scanned_at, approved_at, completed_at, created_at FROM device_pairing_sessions WHERE id = $1`
	
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&session.ID, &session.DeviceID, &session.SessionToken, &session.UserID,
		&session.Purpose, &session.Status, &session.ExpiresAt, &session.ScannedAt,
		&session.ApprovedAt, &session.CompletedAt, &session.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("pairing session not found")
		}
		return nil, err
	}
	return &session, nil
}

func (r *PairingRepository) FindByToken(token string) (*model.DevicePairingSession, error) {
	var session model.DevicePairingSession
	query := `SELECT id, device_id, session_token, user_id, purpose, status, expires_at, scanned_at, approved_at, completed_at, created_at FROM device_pairing_sessions WHERE session_token = $1`
	
	err := r.db.QueryRow(context.Background(), query, token).Scan(
		&session.ID, &session.DeviceID, &session.SessionToken, &session.UserID,
		&session.Purpose, &session.Status, &session.ExpiresAt, &session.ScannedAt,
		&session.ApprovedAt, &session.CompletedAt, &session.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("pairing session not found")
		}
		return nil, err
	}
	return &session, nil
}

func (r *PairingRepository) Update(session *model.DevicePairingSession) error {
	query := `
		UPDATE device_pairing_sessions 
		SET device_id = $1, session_token = $2, user_id = $3, purpose = $4, status = $5, expires_at = $6, scanned_at = $7, approved_at = $8, completed_at = $9
		WHERE id = $10`
	
	commandTag, err := r.db.Exec(context.Background(), query,
		session.DeviceID, session.SessionToken, session.UserID, session.Purpose,
		session.Status, session.ExpiresAt, session.ScannedAt, session.ApprovedAt, session.CompletedAt, session.ID,
	)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("no rows updated")
	}
	return nil
}