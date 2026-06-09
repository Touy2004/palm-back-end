package handler

import (
	"github.com/Touy2004/palm-back-end/internal/service"
	jwtpkg "github.com/Touy2004/palm-back-end/pkg/jwt"
	"github.com/Touy2004/palm-back-end/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var input service.RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Please fill out all required fields correctly.", err.Error())
	}

	user, err := h.authService.Register(input)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "We couldn't register your account. Please check your details and try again.", err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "User registered successfully", fiber.Map{
		"user": fiber.Map{
			"id":            user.ID,
			"employee_code": user.EmployeeCode,
			"full_name":     user.FullName,
			"email":         user.Email,
			"department":    user.Department,
			"phone":         user.Phone,
			"role":          user.Role,
			"status":        user.Status,
		},
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input service.LoginInput
	if err := c.BodyParser(&input); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Please fill out all required fields correctly.", err.Error())
	}

	user, accessToken, refreshToken, err := h.authService.Login(input)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "Your phone number or password is incorrect.", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Login successful", fiber.Map{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&input); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Please provide your refresh token.", err.Error())
	}

	accessToken, newRefreshToken, err := h.authService.RefreshToken(input.RefreshToken)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "Your session has expired. Please log in again.", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Token refreshed successfully", fiber.Map{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwtpkg.Claims)

	user, err := h.authService.GetProfile(claims.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "We couldn't find your user profile.", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Profile retrieved successfully", fiber.Map{
		"user": fiber.Map{
			"id":                 user.ID,
			"employee_code":      user.EmployeeCode,
			"full_name":          user.FullName,
			"email":              user.Email,
			"department":         user.Department,
			"phone":              user.Phone,
			"role":               user.Role,
			"status":             user.Status,
			"is_palm_registered": user.IsPalmRegistered,
		},
	})
}

func (h *AuthHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.authService.GetUsers()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "We encountered an issue fetching the users.", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Users retrieved successfully", users)
}
