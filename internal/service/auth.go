package service

import (
	"errors"

	"github.com/Touy2004/palm-back-end/internal/model"
	"github.com/Touy2004/palm-back-end/internal/repository"
	"github.com/Touy2004/palm-back-end/pkg/hash"
	"github.com/Touy2004/palm-back-end/pkg/jwt"
)

type AuthService struct {
	userRepo *repository.UserRepository
	jwt      *jwt.JWT
}

func NewAuthService(userRepo *repository.UserRepository, jwt *jwt.JWT) *AuthService {
	return &AuthService{userRepo: userRepo, jwt: jwt}
}

type RegisterInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}

type LoginInput struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (s *AuthService) Register(input RegisterInput) (*model.User, error) {
	hashed, err := hash.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.Phone,
		Password:  hashed,
		Role:      input.Role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(input LoginInput) (string, error) {
	user, err := s.userRepo.FindByPhone(input.Phone)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !hash.CheckPassword(input.Password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	return s.jwt.GenerateToken(user.ID, user.Phone, user.Role)
}
func (s *AuthService) GetProfile(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *AuthService) GetUsers() ([]model.User, error) {
	return s.userRepo.FindAll()
}
