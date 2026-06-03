package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/Touy2004/palm-back-end/internal/model"
	"github.com/Touy2004/palm-back-end/internal/repository"
)

type PairingService struct {
	pairingRepo *repository.PairingRepository
}

func NewPairingService(pairingRepo *repository.PairingRepository) *PairingService {
	return &PairingService{pairingRepo: pairingRepo}
}

func (s *PairingService) ScanSession(token string) (*model.DevicePairingSession, error) {
	session, err := s.pairingRepo.FindByToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired session token")
	}

	if session.Status != "pending" {
		return nil, errors.New("session is no longer valid")
	}

	if time.Now().After(session.ExpiresAt) {
		session.Status = "expired"
		_ = s.pairingRepo.Update(session)
		return nil, errors.New("session has expired")
	}

	now := time.Now()
	session.Status = "scanned"
	session.ScannedAt = &now

	if err := s.pairingRepo.Update(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *PairingService) ApproveSession(token, userID string) error {
	session, err := s.pairingRepo.FindByToken(token)
	if err != nil {
		return errors.New("invalid session token")
	}

	if session.Status != "scanned" {
		return errors.New("session must be scanned before approval")
	}

	if time.Now().After(session.ExpiresAt) {
		session.Status = "expired"
		_ = s.pairingRepo.Update(session)
		return errors.New("session has expired")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	now := time.Now()
	session.Status = "approved"
	session.UserID = &userUUID
	session.ApprovedAt = &now

	return s.pairingRepo.Update(session)
}