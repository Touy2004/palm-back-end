package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/Touy2004/palm-back-end/internal/model"
	"github.com/Touy2004/palm-back-end/internal/repository"
)

type DeviceService struct {
	deviceRepo  *repository.DeviceRepository
	pairingRepo *repository.PairingRepository
}

func NewDeviceService(deviceRepo *repository.DeviceRepository, pairingRepo *repository.PairingRepository) *DeviceService {
	return &DeviceService{
		deviceRepo:  deviceRepo,
		pairingRepo: pairingRepo,
	}
}

func (s *DeviceService) Heartbeat(deviceCode string) error {
	device, err := s.deviceRepo.FindByCode(deviceCode)
	if err != nil {
		return errors.New("device not found")
	}

	if device.Status != "active" {
		return errors.New("device is not active")
	}

	now := time.Now()
	device.LastSeenAt = &now
	return s.deviceRepo.Update(device)
}

func (s *DeviceService) CreatePairingSession(deviceCode, purpose string) (*model.DevicePairingSession, error) {
	device, err := s.deviceRepo.FindByCode(deviceCode)
	if err != nil {
		return nil, errors.New("device not found")
	}

	if device.Status != "active" {
		return nil, errors.New("device is not active")
	}

	// Generate a random token
	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}
	token := hex.EncodeToString(tokenBytes)

	session := &model.DevicePairingSession{
		DeviceID:     device.ID,
		SessionToken: token,
		Purpose:      purpose,
		Status:       "pending",
		ExpiresAt:    time.Now().Add(5 * time.Minute), // Expires in 5 minutes
	}

	if err := s.pairingRepo.Create(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *DeviceService) GetSessionStatus(sessionID string) (*model.DevicePairingSession, error) {
	session, err := s.pairingRepo.FindByID(sessionID)
	if err != nil {
		return nil, errors.New("session not found")
	}

	if session.Status == "pending" && time.Now().After(session.ExpiresAt) {
		session.Status = "expired"
		_ = s.pairingRepo.Update(session)
	}

	return session, nil
}
