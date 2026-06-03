package repository

import (
	"time"

	"github.com/Touy2004/palm-back-end/internal/model"
	"gorm.io/gorm"
)

type AttendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{db: db}
}

func (r *AttendanceRepository) FindAll(page, limit int) ([]model.AttendanceLog, int64, error) {
	var logs []model.AttendanceLog
	var total int64

	offset := (page - 1) * limit

	r.db.Model(&model.AttendanceLog{}).Count(&total)
	err := r.db.Preload("User").Preload("Device").
		Order("attendance_date desc").
		Offset(offset).Limit(limit).
		Find(&logs).Error

	return logs, total, err
}

func (r *AttendanceRepository) FindByUserID(userID string, page, limit int) ([]model.AttendanceLog, int64, error) {
	var logs []model.AttendanceLog
	var total int64

	offset := (page - 1) * limit

	query := r.db.Model(&model.AttendanceLog{}).Where("user_id = ?", userID)
	query.Count(&total)

	err := query.Preload("User").Preload("Device").
		Order("attendance_date desc").
		Offset(offset).Limit(limit).
		Find(&logs).Error

	return logs, total, err
}

func (r *AttendanceRepository) FindTodayByUserID(userID string) (*model.AttendanceLog, error) {
	var log model.AttendanceLog
	today := time.Now().Truncate(24 * time.Hour)
	err := r.db.Where("user_id = ? AND attendance_date >= ?", userID, today).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *AttendanceRepository) Create(log *model.AttendanceLog) error {
	return r.db.Create(log).Error
}

func (r *AttendanceRepository) Update(log *model.AttendanceLog) error {
	return r.db.Save(log).Error
}