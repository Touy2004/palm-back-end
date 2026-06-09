package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Device struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	DeviceCode   string     `gorm:"type:varchar(100);unique;not null" json:"device_code"`
	DeviceName   string     `gorm:"column:name;type:varchar(150)" json:"device_name"`
	LocationName string     `gorm:"column:location;type:varchar(150)" json:"location_name"`
	Status       string     `gorm:"type:varchar(30);default:'active'" json:"status"`
	LastSeenAt   *time.Time `json:"last_seen_at"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (d *Device) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return
}