package repository

import (
	"github.com/Touy2004/palm-back-end/internal/model"
	"gorm.io/gorm"
)

type AttemptRepository struct {
	db *gorm.DB
}

func NewAttemptRepository(db *gorm.DB) *AttemptRepository {
	return &AttemptRepository{db: db}
}

func (r *AttemptRepository) Create(attempt *model.PalmAuthAttempt) error {
	return r.db.Create(attempt).Error
}
