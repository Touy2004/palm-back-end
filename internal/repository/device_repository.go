package repository

import (
	"context"
	"errors"

	"github.com/Touy2004/palm-back-end/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeviceRepository struct {
	db *pgxpool.Pool
}

func NewDeviceRepository(db *pgxpool.Pool) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) Create(device *model.Device) error {
	query := `
		INSERT INTO devices (id, device_code, name, location, status, last_seen_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
	
	return r.db.QueryRow(context.Background(), query,
		device.ID, device.DeviceCode, device.DeviceName, device.LocationName,
		device.Status, device.LastSeenAt, device.CreatedAt,
	).Scan(&device.ID)
}

func (r *DeviceRepository) FindAll() ([]model.Device, error) {
	var devices []model.Device
	query := `SELECT id, device_code, name, location, status, last_seen_at, created_at FROM devices`
	
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var device model.Device
		if err := rows.Scan(
			&device.ID, &device.DeviceCode, &device.DeviceName, &device.LocationName,
			&device.Status, &device.LastSeenAt, &device.CreatedAt,
		); err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}
	return devices, rows.Err()
}

func (r *DeviceRepository) FindByID(id string) (*model.Device, error) {
	var device model.Device
	query := `SELECT id, device_code, name, location, status, last_seen_at, created_at FROM devices WHERE id = $1`
	
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&device.ID, &device.DeviceCode, &device.DeviceName, &device.LocationName,
		&device.Status, &device.LastSeenAt, &device.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("device not found")
		}
		return nil, err
	}
	return &device, nil
}

func (r *DeviceRepository) FindByCode(code string) (*model.Device, error) {
	var device model.Device
	query := `SELECT id, device_code, name, location, status, last_seen_at, created_at FROM devices WHERE device_code = $1`
	
	err := r.db.QueryRow(context.Background(), query, code).Scan(
		&device.ID, &device.DeviceCode, &device.DeviceName, &device.LocationName,
		&device.Status, &device.LastSeenAt, &device.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("device not found")
		}
		return nil, err
	}
	return &device, nil
}

func (r *DeviceRepository) Update(device *model.Device) error {
	query := `
		UPDATE devices 
		SET device_code = $1, name = $2, location = $3, status = $4, last_seen_at = $5
		WHERE id = $6`
	
	commandTag, err := r.db.Exec(context.Background(), query,
		device.DeviceCode, device.DeviceName, device.LocationName, device.Status, device.LastSeenAt, device.ID,
	)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("no rows updated")
	}
	return nil
}