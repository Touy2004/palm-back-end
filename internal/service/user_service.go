package service

import (
	"errors"

	"github.com/Touy2004/palm-back-end/internal/model"
	"github.com/Touy2004/palm-back-end/internal/repository"
	"github.com/Touy2004/palm-back-end/pkg/hash"
)

type UserService struct {
	userRepo       *repository.UserRepository
	palmRepo       *repository.PalmRepository
	attendanceRepo *repository.AttendanceRepository
}

func NewUserService(userRepo *repository.UserRepository, palmRepo *repository.PalmRepository, attendanceRepo *repository.AttendanceRepository) *UserService {
	return &UserService{
		userRepo:       userRepo,
		palmRepo:       palmRepo,
		attendanceRepo: attendanceRepo,
	}
}

func (s *UserService) GetPalmTemplates(userID string) ([]model.PalmTemplate, error) {
	return s.palmRepo.FindByUserID(userID)
}

func (s *UserService) DeletePalmTemplate(id, userID string) error {
	return s.palmRepo.Delete(id, userID)
}

func (s *UserService) GetAttendanceHistory(userID string, page, limit int) ([]model.AttendanceLog, int64, error) {
	return s.attendanceRepo.FindByUserID(userID, page, limit)
}

func (s *UserService) ChangePassword(userID, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if !hash.CheckPassword(oldPassword, user.PasswordHash) {
		return errors.New("incorrect old password")
	}

	hashed, err := hash.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	user.PasswordHash = hashed
	return s.userRepo.Update(user)
}
