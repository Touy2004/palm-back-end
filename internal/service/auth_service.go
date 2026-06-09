package service

import (
	"errors"

	"strings"

	"github.com/Touy2004/palm-back-end/internal/constant"
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
	FullName     string `json:"full_name"`
	EmployeeCode string `json:"employee_code"`
	Email        string `json:"email"`
	Department   string `json:"department"`
	Phone        string `json:"phone"`
	Password     string `json:"password"`
	Role         string `json:"role"`
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

	role := strings.ToUpper(input.Role)
	if role == "" {
		role = constant.ROLE_EMPLOYEE
	}

	user := &model.User{
		FullName:     input.FullName,
		EmployeeCode: input.EmployeeCode,
		Email:        input.Email,
		Department:   input.Department,
		Phone:        input.Phone,
		PasswordHash: hashed,
		Role:         role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(input LoginInput) (*model.User, string, string, error) {
	user, err := s.userRepo.FindByPhone(input.Phone)
	if err != nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	if !hash.CheckPassword(input.Password, user.PasswordHash) {
		return nil, "", "", errors.New("invalid credentials")
	}

	accessToken, err := s.jwt.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		return nil, "", "", errors.New("failed to generate access token")
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(user.ID.String(), user.Role)
	if err != nil {
		return nil, "", "", errors.New("failed to generate refresh token")
	}

	return user, accessToken, refreshToken, nil
}

func (s *AuthService) RefreshToken(tokenStr string) (string, string, error) {
	claims, err := s.jwt.Parse(tokenStr)
	if err != nil {
		return "", "", errors.New("invalid or expired refresh token")
	}

	accessToken, err := s.jwt.GenerateToken(claims.UserID, claims.Role)
	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(claims.UserID, claims.Role)
	if err != nil {
		return "", "", errors.New("failed to generate new refresh token")
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) GetProfile(userID string) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *AuthService) GetUsers() ([]model.User, error) {
	return s.userRepo.FindAll()
}
