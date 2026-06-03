package repository

import (
	"github.com/Touy2004/palm-back-end/internal/model"
	"gorm.io/gorm"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) Create(device *model.Device) error {
	return r.db.Create(device).Error
}

func (r *DeviceRepository) FindAll() ([]model.Device, error) {
	var devices []model.Device
	err := r.db.Find(&devices).Error
	return devices, err
}

func (r *DeviceRepository) FindByID(id string) (*model.Device, error) {
	var device model.Device
	err := r.db.First(&device, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *DeviceRepository) FindByCode(code string) (*model.Device, error) {
	var device model.Device
	err := r.db.Where("device_code = ?", code).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *DeviceRepository) Update(device *model.Device) error {
	return r.db.Save(device).Error
}