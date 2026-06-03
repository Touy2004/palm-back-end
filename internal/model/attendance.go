package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttendanceLog struct {
	ID               uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID           uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_user_date" json:"user_id"`
	DeviceID         *uuid.UUID `gorm:"type:uuid" json:"device_id"`
	AttendanceDate   time.Time  `gorm:"type:date;not null;uniqueIndex:idx_user_date;index:idx_date" json:"attendance_date"`
	CheckInTime      *time.Time `json:"check_in_time"`
	CheckOutTime     *time.Time `json:"check_out_time"`
	CheckInScore     *float64   `gorm:"type:numeric(6,5)" json:"check_in_score"`
	CheckOutScore    *float64   `gorm:"type:numeric(6,5)" json:"check_out_score"`
	CheckInLiveness  *bool      `json:"check_in_liveness"`
	CheckOutLiveness *bool      `json:"check_out_liveness"`
	Status           string     `gorm:"type:varchar(30);default:'present'" json:"status"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Device *Device `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
}

func (a *AttendanceLog) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}