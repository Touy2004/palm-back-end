package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository struct {
	db *pgxpool.Pool
}

func NewAdminRepository(db *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{db: db}
}

type DashboardSummary struct {
	TotalUsers          int `json:"total_users"`
	TotalDevices        int `json:"total_devices"`
	ActivePalmTemplates int `json:"active_palm_templates"`
	CheckInsToday       int `json:"check_ins_today"`
}

func (r *AdminRepository) GetDashboardSummary() (*DashboardSummary, error) {
	var summary DashboardSummary

	// Total Users
	err := r.db.QueryRow(context.Background(), `SELECT count(*) FROM users`).Scan(&summary.TotalUsers)
	if err != nil {
		return nil, err
	}

	// Total Devices
	err = r.db.QueryRow(context.Background(), `SELECT count(*) FROM devices`).Scan(&summary.TotalDevices)
	if err != nil {
		return nil, err
	}

	// Active Palm Templates
	err = r.db.QueryRow(context.Background(), `SELECT count(*) FROM palm_templates WHERE status = 'active'`).Scan(&summary.ActivePalmTemplates)
	if err != nil {
		return nil, err
	}

	// Check-ins Today
	err = r.db.QueryRow(context.Background(), `
		SELECT count(*) 
		FROM attendance_logs 
		WHERE DATE(attendance_date) = CURRENT_DATE
	`).Scan(&summary.CheckInsToday)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}
