package handler

import (
	"strconv"

	"github.com/Touy2004/palm-back-end/internal/service"
	jwtpkg "github.com/Touy2004/palm-back-end/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService    *service.UserService
	pairingService *service.PairingService
}

func NewUserHandler(userService *service.UserService, pairingService *service.PairingService) *UserHandler {
	return &UserHandler{
		userService:    userService,
		pairingService: pairingService,
	}
}

// Pairing
func (h *UserHandler) ScanPairingQR(c *fiber.Ctx) error {
	var input struct {
		SessionToken string `json:"session_token"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	session, err := h.pairingService.ScanSession(input.SessionToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pairing session scanned",
		"device":  session.Device,
		"purpose": session.Purpose,
	})
}

func (h *UserHandler) ApprovePairingQR(c *fiber.Ctx) error {
	var input struct {
		SessionToken string `json:"session_token"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	claims := c.Locals("user").(*jwtpkg.Claims)

	err := h.pairingService.ApproveSession(input.SessionToken, claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Enrollment approved. Please place your palm on the device.",
	})
}

// Palm Templates
func (h *UserHandler) GetPalmTemplates(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwtpkg.Claims)

	templates, err := h.userService.GetPalmTemplates(claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch palm templates"})
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"templates": templates,
	})
}

func (h *UserHandler) DeletePalmTemplate(c *fiber.Ctx) error {
	id := c.Params("id")
	claims := c.Locals("user").(*jwtpkg.Claims)

	err := h.userService.DeletePalmTemplate(id, claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete palm template"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Palm template deleted",
	})
}

// Attendance
func (h *UserHandler) GetMyAttendance(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwtpkg.Claims)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	logs, total, err := h.userService.GetAttendanceHistory(claims.UserID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch attendance history"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    logs,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	var input struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if len(input.NewPassword) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "new password must be at least 6 characters"})
	}

	claims := c.Locals("user").(*jwtpkg.Claims)

	if err := h.userService.ChangePassword(claims.UserID, input.OldPassword, input.NewPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Password changed successfully",
	})
}
