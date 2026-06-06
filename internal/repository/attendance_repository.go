package repository

import (
	"context"
	"errors"
	"fmt"
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

func (r *AttendanceRepository) FindAll(page, limit int, startDate, endDate string) ([]model.AttendanceLog, int64, error) {
	var logs []model.AttendanceLog
	var total int64
	
	countQuery := `SELECT count(*) FROM attendance_logs WHERE 1=1`
	var countArgs []interface{}
	argIdx := 1

	if startDate != "" {
		countQuery += fmt.Sprintf(" AND attendance_date >= $%d", argIdx)
		countArgs = append(countArgs, startDate)
		argIdx++
	}
	if endDate != "" {
		countQuery += fmt.Sprintf(" AND attendance_date <= $%d", argIdx)
		countArgs = append(countArgs, endDate)
		argIdx++
	}

	err := r.db.QueryRow(context.Background(), countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := `
		SELECT a.id, a.user_id, a.device_id, a.attendance_date, a.check_in_time, a.check_out_time, 
		       a.check_in_score, a.check_out_score, a.check_in_liveness, a.check_out_liveness, a.status, a.created_at,
		       d.name, d.device_code
		FROM attendance_logs a
		LEFT JOIN devices d ON a.device_id = d.id
		WHERE 1=1`
	
	var args []interface{}
	argIdx = 1

	if startDate != "" {
		query += fmt.Sprintf(" AND a.attendance_date >= $%d", argIdx)
		args = append(args, startDate)
		argIdx++
	}
	if endDate != "" {
		query += fmt.Sprintf(" AND a.attendance_date <= $%d", argIdx)
		args = append(args, endDate)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY a.attendance_date DESC, a.created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)
	
	rows, err := r.db.Query(context.Background(), query, args...)
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
			&log.DeviceName, &log.DeviceCode,
		); err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}
	return logs, total, rows.Err()
}

func (r *AttendanceRepository) FindByUserID(userID string, page, limit int, startDate, endDate string) ([]model.AttendanceLog, int64, error) {
	var logs []model.AttendanceLog
	var total int64
	
	countQuery := `SELECT count(*) FROM attendance_logs WHERE user_id = $1`
	countArgs := []interface{}{userID}
	argIdx := 2

	if startDate != "" {
		countQuery += fmt.Sprintf(" AND attendance_date >= $%d", argIdx)
		countArgs = append(countArgs, startDate)
		argIdx++
	}
	if endDate != "" {
		countQuery += fmt.Sprintf(" AND attendance_date <= $%d", argIdx)
		countArgs = append(countArgs, endDate)
		argIdx++
	}

	err := r.db.QueryRow(context.Background(), countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := `
		SELECT a.id, a.user_id, a.device_id, a.attendance_date, a.check_in_time, a.check_out_time, 
		       a.check_in_score, a.check_out_score, a.check_in_liveness, a.check_out_liveness, a.status, a.created_at,
		       d.name, d.device_code
		FROM attendance_logs a
		LEFT JOIN devices d ON a.device_id = d.id
		WHERE a.user_id = $1`
	
	args := []interface{}{userID}
	argIdx = 2

	if startDate != "" {
		query += fmt.Sprintf(" AND a.attendance_date >= $%d", argIdx)
		args = append(args, startDate)
		argIdx++
	}
	if endDate != "" {
		query += fmt.Sprintf(" AND a.attendance_date <= $%d", argIdx)
		args = append(args, endDate)
		argIdx++
	}
	
	query += fmt.Sprintf(" ORDER BY a.attendance_date DESC, a.created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)
	
	rows, err := r.db.Query(context.Background(), query, args...)
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
			&log.DeviceName, &log.DeviceCode,
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
	
	query := `
		SELECT a.id, a.user_id, a.device_id, a.attendance_date, a.check_in_time, a.check_out_time, 
		       a.check_in_score, a.check_out_score, a.check_in_liveness, a.check_out_liveness, a.status, a.created_at,
		       d.name, d.device_code
		FROM attendance_logs a
		LEFT JOIN devices d ON a.device_id = d.id
		WHERE a.user_id = $1 AND a.attendance_date >= $2`
	
	err := r.db.QueryRow(context.Background(), query, userID, today).Scan(
		&log.ID, &log.UserID, &log.DeviceID, &log.AttendanceDate,
		&log.CheckInTime, &log.CheckOutTime, &log.CheckInScore, &log.CheckOutScore,
		&log.CheckInLiveness, &log.CheckOutLiveness, &log.Status, &log.CreatedAt,
		&log.DeviceName, &log.DeviceCode,
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

func (r *AttendanceRepository) FindAllByDateRange(startDate, endDate string) ([]model.AttendanceLog, error) {
	query := `
		SELECT a.id, a.user_id, a.device_id, a.attendance_date, a.check_in_time, a.check_out_time, 
		       a.check_in_score, a.check_out_score, a.check_in_liveness, a.check_out_liveness, a.status, a.created_at,
		       d.name, d.device_code
		FROM attendance_logs a
		LEFT JOIN devices d ON a.device_id = d.id
		WHERE a.attendance_date >= $1 AND a.attendance_date <= $2
	`
	
	rows, err := r.db.Query(context.Background(), query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []model.AttendanceLog
	for rows.Next() {
		var log model.AttendanceLog
		var deviceID *uuid.UUID
		var checkInTime, checkOutTime *time.Time
		var checkInScore, checkOutScore *float64
		var deviceName, deviceCode *string

		err := rows.Scan(
			&log.ID, &log.UserID, &deviceID, &log.AttendanceDate, &checkInTime, &checkOutTime,
			&checkInScore, &checkOutScore, &log.CheckInLiveness, &log.CheckOutLiveness, &log.Status, &log.CreatedAt,
			&deviceName, &deviceCode,
		)
		if err != nil {
			return nil, err
		}

		if deviceID != nil {
			log.DeviceID = deviceID
		}
		if checkInTime != nil {
			log.CheckInTime = checkInTime
		}
		if checkOutTime != nil {
			log.CheckOutTime = checkOutTime
		}
		if checkInScore != nil {
			log.CheckInScore = checkInScore
		}
		if checkOutScore != nil {
			log.CheckOutScore = checkOutScore
		}
		if deviceName != nil {
			log.DeviceName = deviceName
		}
		if deviceCode != nil {
			log.DeviceCode = deviceCode
		}

		logs = append(logs, log)
	}

	return logs, nil
}