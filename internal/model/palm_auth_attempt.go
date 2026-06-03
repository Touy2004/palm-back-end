package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PalmAuthAttempt struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID         *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	DeviceID       *uuid.UUID `gorm:"type:uuid" json:"device_id"`
	TemplateID     *uuid.UUID `gorm:"type:uuid" json:"template_id"`
	Action         string     `gorm:"type:varchar(50);not null" json:"action"`
	Score          *float64   `gorm:"type:numeric(6,5)" json:"score"`
	Threshold      *float64   `gorm:"type:numeric(5,4)" json:"threshold"`
	LivenessPassed bool       `gorm:"default:false" json:"liveness_passed"`
	QualityScore   *float64   `gorm:"type:numeric(6,5)" json:"quality_score"`
	ThermalMin     *float64   `gorm:"type:numeric(6,2)" json:"thermal_min"`
	ThermalMax     *float64   `gorm:"type:numeric(6,2)" json:"thermal_max"`
	ThermalAvg     *float64   `gorm:"type:numeric(6,2)" json:"thermal_avg"`
	Result         string     `gorm:"type:varchar(30);not null" json:"result"`
	FailureReason  string     `gorm:"type:text" json:"failure_reason"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (a *PalmAuthAttempt) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}
