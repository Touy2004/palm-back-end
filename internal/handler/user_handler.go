package handler

import (
	"strconv"

	"github.com/Touy2004/palm-back-end/internal/model"
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
		return response.Error(c, fiber.StatusBadRequest, "Please make sure your QR code scan is valid.", err.Error())
	}

	session, err := h.pairingService.ScanSession(input.SessionToken)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "This QR code is invalid or has already expired.", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Pairing session scanned", fiber.Map{
		"device":  session.Device,
		"purpose": session.Purpose,
	})
}

func (h *UserHandler) ApprovePairingQR(c *fiber.Ctx) error {
	var input struct {
		SessionToken string `json:"session_token"`
		HandSide     string `json:"hand_side"`
	}
	if err := c.BodyParser(&input); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Please provide valid approval data.", err.Error())
	}

	if input.HandSide == "" {
		return response.Error(c, fiber.StatusBadRequest, "Please select whether you want to use your left or right hand.", nil)
	}

	claims := c.Locals("user").(*jwtpkg.Claims)

	err := h.pairingService.ApproveSession(input.SessionToken, claims.UserID, input.HandSide)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "We couldn't approve the scanner. Please try scanning the QR code again.", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Enrollment approved. Please place your palm on the device.", nil)
}

// Palm Templates
func (h *UserHandler) GetPalmTemplates(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwtpkg.Claims)

	templates, err := h.userService.GetPalmTemplates(claims.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "We encountered an issue fetching your palm templates.", err.Error())
	}

	if templates == nil {
		templates = make([]model.PalmTemplate, 0)
	}

	return response.Success(c, fiber.StatusOK, "Palm templates retrieved successfully", templates)
}

func (h *UserHandler) DeletePalmTemplate(c *fiber.Ctx) error {
	id := c.Params("id")
	claims := c.Locals("user").(*jwtpkg.Claims)

	err := h.userService.DeletePalmTemplate(id, claims.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "We couldn't delete your palm template.", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Palm template deleted successfully", nil)
}

// Attendance
func (h *UserHandler) GetMyAttendance(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwtpkg.Claims)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	logs, total, err := h.userService.GetAttendanceHistory(claims.UserID, page, limit, startDate, endDate)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "We encountered an issue fetching your attendance history.", err.Error())
	}

	if logs == nil {
		logs = make(model.AttendanceLogs, 0)
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
		return response.Error(c, fiber.StatusBadRequest, "Please provide your old and new password correctly.", err.Error())
	}

	if len(input.NewPassword) < 6 {
		return response.Error(c, fiber.StatusBadRequest, "Your new password must be at least 6 characters long.", "new password must be at least 6 characters")
	}

	claims := c.Locals("user").(*jwtpkg.Claims)

	if err := h.userService.ChangePassword(claims.UserID, input.OldPassword, input.NewPassword); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Your old password is incorrect, or we couldn't change it.", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Password changed successfully", nil)
}
