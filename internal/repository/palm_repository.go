package repository

import (
	"github.com/Touy2004/palm-back-end/internal/model"
	"gorm.io/gorm"
)

type PalmRepository struct {
	db *gorm.DB
}

func NewPalmRepository(db *gorm.DB) *PalmRepository {
	return &PalmRepository{db: db}
}

func (r *PalmRepository) Create(template *model.PalmTemplate) error {
	return r.db.Create(template).Error
}

func (r *PalmRepository) FindByUserID(userID string) ([]model.PalmTemplate, error) {
	var templates []model.PalmTemplate
	err := r.db.Where("user_id = ? AND status = 'active'", userID).Find(&templates).Error
	return templates, err
}

func (r *PalmRepository) Delete(id, userID string) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.PalmTemplate{}).Error
}

func (r *PalmRepository) FindAllActive() ([]model.PalmTemplate, error) {
	var templates []model.PalmTemplate
	err := r.db.Where("status = 'active'").Find(&templates).Error
	return templates, err
}