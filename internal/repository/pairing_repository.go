package repository

import (
	"github.com/Touy2004/palm-back-end/internal/model"
	"gorm.io/gorm"
)

type PairingRepository struct {
	db *gorm.DB
}

func NewPairingRepository(db *gorm.DB) *PairingRepository {
	return &PairingRepository{db: db}
}

func (r *PairingRepository) Create(session *model.DevicePairingSession) error {
	return r.db.Create(session).Error
}

func (r *PairingRepository) FindByID(id string) (*model.DevicePairingSession, error) {
	var session model.DevicePairingSession
	err := r.db.Preload("User").Preload("Device").First(&session, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *PairingRepository) FindByToken(token string) (*model.DevicePairingSession, error) {
	var session model.DevicePairingSession
	err := r.db.Preload("Device").Where("session_token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *PairingRepository) Update(session *model.DevicePairingSession) error {
	return r.db.Save(session).Error
}