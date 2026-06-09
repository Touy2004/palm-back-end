package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DevicePairingSession struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	DeviceID     uuid.UUID  `gorm:"type:uuid;not null" json:"device_id"`
	SessionToken string     `gorm:"type:text;unique;not null" json:"session_token"`
	UserID       *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	HandSide     *string    `gorm:"type:varchar(10)" json:"hand_side"`
	Purpose      string     `gorm:"type:varchar(50);not null" json:"purpose"`
	Status       string     `gorm:"type:varchar(30);default:'pending'" json:"status"`
	ExpiresAt    time.Time  `gorm:"not null" json:"expires_at"`
	ScannedAt    *time.Time `json:"scanned_at"`
	ApprovedAt   *time.Time `json:"approved_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`

	Device *Device `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (s *DevicePairingSession) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
