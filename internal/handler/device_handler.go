package handler

import (
	"github.com/Touy2004/palm-back-end/internal/service"
	"github.com/gofiber/fiber/v2"
)

type DeviceHandler struct {
	deviceSvc     *service.DeviceService
	palmSvc       *service.PalmService
	attendanceSvc *service.AttendanceService
}

func NewDeviceHandler(
	deviceSvc *service.DeviceService,
	palmSvc *service.PalmService,
	attendanceSvc *service.AttendanceService,
) *DeviceHandler {
	return &DeviceHandler{
		deviceSvc:     deviceSvc,
		palmSvc:       palmSvc,
		attendanceSvc: attendanceSvc,
	}
}

func (h *DeviceHandler) Heartbeat(c *fiber.Ctx) error {
	var input struct {
		DeviceCode string `json:"device_code"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if err := h.deviceSvc.Heartbeat(input.DeviceCode); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

func (h *DeviceHandler) CreatePairingSession(c *fiber.Ctx) error {
	var input struct {
		DeviceCode string `json:"device_code"`
		Purpose    string `json:"purpose"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	session, err := h.deviceSvc.CreatePairingSession(input.DeviceCode, input.Purpose)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"session_id":    session.ID,
		"session_token": session.SessionToken,
		"expires_at":    session.ExpiresAt,
	})
}

func (h *DeviceHandler) GetSessionStatus(c *fiber.Ctx) error {
	sessionID := c.Params("session_id")
	session, err := h.deviceSvc.GetSessionStatus(sessionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"status":  session.Status,
	})
}

func (h *DeviceHandler) EnrollPalm(c *fiber.Ctx) error {
	var payload struct {
		DeviceCode     string      `json:"device_code"`
		SessionToken   string      `json:"session_token"`
		HandSide       string      `json:"hand_side"`
		ModelVersion   string      `json:"model_version"`
		EmbeddingDim   int         `json:"embedding_dim"`
		Embeddings     [][]float32 `json:"embeddings"`
		LivenessPassed bool        `json:"liveness_passed"`
		QualityScore   float64     `json:"quality_score"`
		ThermalMin     float64     `json:"thermal_min"`
		ThermalMax     float64     `json:"thermal_max"`
		ThermalAvg     float64     `json:"thermal_avg"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	enrollInput := service.EnrollInput{
		SessionToken:   payload.SessionToken,
		DeviceCode:     payload.DeviceCode, // Passed just in case, but session has it
		HandSide:       payload.HandSide,
		ModelVersion:   payload.ModelVersion,
		EmbeddingDim:   payload.EmbeddingDim,
		Embeddings:     payload.Embeddings,
		LivenessPassed: payload.LivenessPassed,
		QualityScore:   payload.QualityScore,
		ThermalMin:     payload.ThermalMin,
		ThermalMax:     payload.ThermalMax,
		ThermalAvg:     payload.ThermalAvg,
	}

	template, err := h.palmSvc.EnrollPalm(enrollInput)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success":     true,
		"template_id": template.ID,
		"message":     "Palm enrolled successfully",
	})
}

func (h *DeviceHandler) ProcessAttendance(c *fiber.Ctx) error {
	var payload service.ProcessAttendanceInput
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	result, err := h.attendanceSvc.ProcessPalmAttendance(payload)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"action":  result.Action,
		"user": fiber.Map{
			"id":        result.UserID,
			"full_name": result.FullName,
		},
		"message": result.Message,
	})
}

func (h *DeviceHandler) IdentifyPalm(c *fiber.Ctx) error {
	var payload service.ProcessAttendanceInput
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	result, err := h.attendanceSvc.IdentifyPalm(payload)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user": fiber.Map{
			"id":        result.UserID,
			"full_name": result.FullName,
		},
		"message": result.Message,
	})
}