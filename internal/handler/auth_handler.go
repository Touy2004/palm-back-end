package handler

import (
	"github.com/Touy2004/palm-back-end/internal/service"
	jwtpkg "github.com/Touy2004/palm-back-end/pkg/jwt"

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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	user, err := h.authService.Register(input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	user, accessToken, refreshToken, err := h.authService.Login(input)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success":       true,
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	accessToken, newRefreshToken, err := h.authService.RefreshToken(input.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwtpkg.Claims)

	user, err := h.authService.GetProfile(claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return c.JSON(fiber.Map{
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

func (h *AuthHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.authService.GetUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not fetch users",
		})
	}

	return c.JSON(fiber.Map{
		"users": users,
	})
}
