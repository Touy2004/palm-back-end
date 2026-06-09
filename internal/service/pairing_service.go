package service

import (
	"errors"
	"time"

	"github.com/Touy2004/palm-back-end/internal/model"
	"github.com/Touy2004/palm-back-end/internal/repository"
	"github.com/google/uuid"
)

type PairingService struct {
	pairingRepo *repository.PairingRepository
	palmRepo    *repository.PalmRepository
}

func NewPairingService(pairingRepo *repository.PairingRepository, palmRepo *repository.PalmRepository) *PairingService {
	return &PairingService{
		pairingRepo: pairingRepo,
		palmRepo:    palmRepo,
	}
}

func (s *PairingService) ScanSession(token string) (*model.DevicePairingSession, error) {
	session, err := s.pairingRepo.FindByToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired session token: " + err.Error())
	}

	if session.Status != "pending" {
		return nil, errors.New("session is no longer valid")
	}

	if time.Now().UTC().After(session.ExpiresAt) {
		session.Status = "expired"
		_ = s.pairingRepo.Update(session)
		return nil, errors.New("session has expired")
	}

	now := time.Now().UTC()
	session.Status = "scanned"
	session.ScannedAt = &now

	if err := s.pairingRepo.Update(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *PairingService) ApproveSession(token, userID, handSide string) error {
	if handSide != "left" && handSide != "right" {
		return errors.New("invalid hand side, must be 'left' or 'right'")
	}

	session, err := s.pairingRepo.FindByToken(token)
	if err != nil {
		return errors.New("invalid session token")
	}

	if session.Status != "scanned" {
		return errors.New("session must be scanned before approval")
	}

	if time.Now().UTC().After(session.ExpiresAt) {
		session.Status = "expired"
		_ = s.pairingRepo.Update(session)
		return errors.New("session has expired")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Check if user already has an active template for this hand
	existingTemplates, err := s.palmRepo.FindByUserID(userID)
	if err == nil {
		for _, t := range existingTemplates {
			if t.HandSide == handSide {
				return errors.New("you already have a registered template for your " + handSide + " hand. Please delete it first")
			}
		}
	}

	now := time.Now().UTC()
	session.Status = "approved"
	session.UserID = &userUUID
	session.HandSide = &handSide
	session.ApprovedAt = &now

	return s.pairingRepo.Update(session)
}
