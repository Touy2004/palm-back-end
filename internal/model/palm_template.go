package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PalmTemplate struct {
	ID                 uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID             uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	HandSide           string     `gorm:"type:varchar(10);not null" json:"hand_side"`
	TemplateEncrypted  []byte     `gorm:"type:bytea;not null" json:"-"`
	TemplateNonce      []byte     `gorm:"type:bytea;not null" json:"-"`
	EmbeddingDim       int        `gorm:"not null;default:128" json:"embedding_dim"`
	ModelVersion       string     `gorm:"type:varchar(100);not null" json:"model_version"`
	Threshold          float64    `gorm:"type:numeric(5,4);default:0.8200" json:"threshold"`
	Status             string     `gorm:"type:varchar(30);default:'active'" json:"status"`
	RegisteredDeviceID *uuid.UUID `gorm:"type:uuid" json:"registered_device_id"`
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	RevokedAt          *time.Time `json:"revoked_at"`

	User             *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	RegisteredDevice *Device `gorm:"foreignKey:RegisteredDeviceID" json:"registered_device,omitempty"`
}

func (t *PalmTemplate) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}
