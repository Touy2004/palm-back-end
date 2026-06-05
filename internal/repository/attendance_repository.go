package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Touy2004/palm-back-end/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AttendanceRepository struct {
	db *pgxpool.Pool
}

func NewAttendanceRepository(db *pgxpool.Pool) *AttendanceRepository {
	return &AttendanceRepository{db: db}
}

func (r *AttendanceRepository) FindAll(page, limit int) ([]model.AttendanceLog, int64, error) {
	var logs []model.AttendanceLog
	var total int64
	
	countQuery := `SELECT count(*) FROM attendance_logs`
	err := r.db.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := `SELECT id, user_id, device_id, attendance_date, check_in_time, check_out_time, check_in_score, check_out_score, check_in_liveness, check_out_liveness, status, created_at FROM attendance_logs ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var log model.AttendanceLog
		if err := rows.Scan(
			&log.ID, &log.UserID, &log.DeviceID, &log.AttendanceDate,
			&log.CheckInTime, &log.CheckOutTime, &log.CheckInScore, &log.CheckOutScore,
			&log.CheckInLiveness, &log.CheckOutLiveness, &log.Status, &log.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}
	return logs, total, rows.Err()
}

func (r *AttendanceRepository) FindByUserID(userID string, page, limit int) ([]model.AttendanceLog, int64, error) {
	var logs []model.AttendanceLog
	var total int64
	
	countQuery := `SELECT count(*) FROM attendance_logs WHERE user_id = $1`
	err := r.db.QueryRow(context.Background(), countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := `SELECT id, user_id, device_id, attendance_date, check_in_time, check_out_time, check_in_score, check_out_score, check_in_liveness, check_out_liveness, status, created_at FROM attendance_logs WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	
	rows, err := r.db.Query(context.Background(), query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var log model.AttendanceLog
		if err := rows.Scan(
			&log.ID, &log.UserID, &log.DeviceID, &log.AttendanceDate,
			&log.CheckInTime, &log.CheckOutTime, &log.CheckInScore, &log.CheckOutScore,
			&log.CheckInLiveness, &log.CheckOutLiveness, &log.Status, &log.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}
	return logs, total, rows.Err()
}

func (r *AttendanceRepository) FindTodayByUserID(userID string) (*model.AttendanceLog, error) {
	var log model.AttendanceLog
	today := time.Now().Truncate(24 * time.Hour)
	
	query := `SELECT id, user_id, device_id, attendance_date, check_in_time, check_out_time, check_in_score, check_out_score, check_in_liveness, check_out_liveness, status, created_at FROM attendance_logs WHERE user_id = $1 AND attendance_date >= $2`
	
	err := r.db.QueryRow(context.Background(), query, userID, today).Scan(
		&log.ID, &log.UserID, &log.DeviceID, &log.AttendanceDate,
		&log.CheckInTime, &log.CheckOutTime, &log.CheckInScore, &log.CheckOutScore,
		&log.CheckInLiveness, &log.CheckOutLiveness, &log.Status, &log.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("attendance log not found")
		}
		return nil, err
	}
	return &log, nil
}

func (r *AttendanceRepository) Create(log *model.AttendanceLog) error {
	if log.ID == uuid.Nil {
		log.ID = uuid.New()
	}
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO attendance_logs (id, user_id, device_id, attendance_date, check_in_time, check_out_time, check_in_score, check_out_score, check_in_liveness, check_out_liveness, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at`
	
	return r.db.QueryRow(context.Background(), query,
		log.ID, log.UserID, log.DeviceID, log.AttendanceDate,
		log.CheckInTime, log.CheckOutTime, log.CheckInScore, log.CheckOutScore,
		log.CheckInLiveness, log.CheckOutLiveness, log.Status, log.CreatedAt,
	).Scan(&log.ID, &log.CreatedAt)
}

func (r *AttendanceRepository) Update(log *model.AttendanceLog) error {
	query := `
		UPDATE attendance_logs 
		SET device_id = $1, check_in_time = $2, check_out_time = $3, check_in_score = $4, check_out_score = $5, check_in_liveness = $6, check_out_liveness = $7, status = $8
		WHERE id = $9`
	
	commandTag, err := r.db.Exec(context.Background(), query,
		log.DeviceID, log.CheckInTime, log.CheckOutTime, log.CheckInScore, log.CheckOutScore,
		log.CheckInLiveness, log.CheckOutLiveness, log.Status, log.ID,
	)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("no rows updated")
	}
	return nil
}