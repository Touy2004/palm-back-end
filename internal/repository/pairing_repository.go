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
		INSERT INTO device_pairing_sessions (id, device_id, session_token, user_id, hand_side, purpose, status, expires_at, scanned_at, approved_at, completed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at`
	
	return r.db.QueryRow(context.Background(), query,
		session.ID, session.DeviceID, session.SessionToken, session.UserID, session.HandSide,
		session.Purpose, session.Status, session.ExpiresAt, session.ScannedAt,
		session.ApprovedAt, session.CompletedAt, session.CreatedAt,
	).Scan(&session.ID, &session.CreatedAt)
}

func (r *PairingRepository) FindByID(id string) (*model.DevicePairingSession, error) {
	var session model.DevicePairingSession
	query := `SELECT id, device_id, session_token, user_id, hand_side, purpose, status, expires_at, scanned_at, approved_at, completed_at, created_at FROM device_pairing_sessions WHERE id = $1`
	
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&session.ID, &session.DeviceID, &session.SessionToken, &session.UserID, &session.HandSide,
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
	session.Device = &model.Device{}
	query := `
		SELECT s.id, s.device_id, s.session_token, s.user_id, s.hand_side, s.purpose, s.status, s.expires_at, s.scanned_at, s.approved_at, s.completed_at, s.created_at,
		       d.id, d.device_code, COALESCE(d.device_name, ''), COALESCE(d.location_name, '')
		FROM device_pairing_sessions s
		LEFT JOIN devices d ON s.device_id = d.id
		WHERE s.session_token = $1
	`
	
	err := r.db.QueryRow(context.Background(), query, token).Scan(
		&session.ID, &session.DeviceID, &session.SessionToken, &session.UserID, &session.HandSide,
		&session.Purpose, &session.Status, &session.ExpiresAt, &session.ScannedAt,
		&session.ApprovedAt, &session.CompletedAt, &session.CreatedAt,
		&session.Device.ID, &session.Device.DeviceCode, &session.Device.DeviceName, &session.Device.LocationName,
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
		SET device_id = $1, session_token = $2, user_id = $3, hand_side = $4, purpose = $5, status = $6, expires_at = $7, scanned_at = $8, approved_at = $9, completed_at = $10
		WHERE id = $11`
	
	commandTag, err := r.db.Exec(context.Background(), query,
		session.DeviceID, session.SessionToken, session.UserID, session.HandSide, session.Purpose,
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