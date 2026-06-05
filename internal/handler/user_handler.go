package handler

import (
	"strconv"

	"github.com/Touy2004/palm-back-end/internal/service"
	jwtpkg "github.com/Touy2004/palm-back-end/pkg/jwt"
	"github.com/Touy2004/palm-back-end/pkg/response"
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
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	session, err := h.pairingService.ScanSession(input.SessionToken)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Failed to scan session", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Pairing session scanned", fiber.Map{
		"device":  session.Device,
		"purpose": session.Purpose,
	})
}

func (h *UserHandler) ApprovePairingQR(c *fiber.Ctx) error {
	var input struct {
		SessionToken string `json:"session_token"`
	}
	if err := c.BodyParser(&input); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	claims := c.Locals("user").(*jwtpkg.Claims)

	err := h.pairingService.ApproveSession(input.SessionToken, claims.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Failed to approve session", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Enrollment approved. Please place your palm on the device.", nil)
}

// Palm Templates
func (h *UserHandler) GetPalmTemplates(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwtpkg.Claims)

	templates, err := h.userService.GetPalmTemplates(claims.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch palm templates", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Palm templates retrieved successfully", templates)
}

func (h *UserHandler) DeletePalmTemplate(c *fiber.Ctx) error {
	id := c.Params("id")
	claims := c.Locals("user").(*jwtpkg.Claims)

	err := h.userService.DeletePalmTemplate(id, claims.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete palm template", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Palm template deleted successfully", nil)
}

// Attendance
func (h *UserHandler) GetMyAttendance(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwtpkg.Claims)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	logs, total, err := h.userService.GetAttendanceHistory(claims.UserID, page, limit)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch attendance history", err.Error())
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "Attendance history retrieved successfully", logs, fiber.Map{
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
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if len(input.NewPassword) < 6 {
		return response.Error(c, fiber.StatusBadRequest, "Invalid password", "new password must be at least 6 characters")
	}

	claims := c.Locals("user").(*jwtpkg.Claims)

	if err := h.userService.ChangePassword(claims.UserID, input.OldPassword, input.NewPassword); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Failed to change password", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Password changed successfully", nil)
}
