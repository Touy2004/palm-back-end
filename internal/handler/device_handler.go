package handler

import (
	"github.com/Touy2004/palm-back-end/internal/service"
	"github.com/Touy2004/palm-back-end/pkg/response"
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
		return response.Error(c, fiber.StatusBadRequest, "Please provide a valid device code.", err.Error())
	}

	if err := h.deviceSvc.Heartbeat(input.DeviceCode); err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "The device is unauthorized to perform this action.", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Heartbeat successful", nil)
}

func (h *DeviceHandler) CreatePairingSession(c *fiber.Ctx) error {
	var input struct {
		DeviceCode string `json:"device_code"`
		Purpose    string `json:"purpose"`
	}
	if err := c.BodyParser(&input); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Please provide valid pairing data.", err.Error())
	}

	session, err := h.deviceSvc.CreatePairingSession(input.DeviceCode, input.Purpose)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to initiate pairing. Please check the device code.", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Pairing session created successfully", fiber.Map{
		"session_id":    session.ID,
		"session_token": session.SessionToken,
		"expires_at":    session.ExpiresAt,
	})
}

func (h *DeviceHandler) GetSessionStatus(c *fiber.Ctx) error {
	sessionID := c.Params("session_id")
	session, err := h.deviceSvc.GetSessionStatus(sessionID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "The pairing session could not be found or has expired.", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Session status retrieved", fiber.Map{
		"status": session.Status,
	})
}

func (h *DeviceHandler) EnrollPalm(c *fiber.Ctx) error {
	var payload struct {
		DeviceCode     string      `json:"device_code"`
		SessionToken   string      `json:"session_token"`
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
		return response.Error(c, fiber.StatusBadRequest, "Please provide a valid scan payload.", err.Error())
	}

	enrollInput := service.EnrollInput{
		SessionToken:   payload.SessionToken,
		DeviceCode:     payload.DeviceCode, // Passed just in case, but session has it
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
		return response.Error(c, fiber.StatusBadRequest, "Palm enrollment failed. Please ensure your hand is placed flat and try again.", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Palm enrolled successfully", fiber.Map{
		"template_id": template.ID,
	})
}

func (h *DeviceHandler) ProcessAttendance(c *fiber.Ctx) error {
	var payload struct {
		DeviceCode     string      `json:"device_code"`
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
		return response.Error(c, fiber.StatusBadRequest, "Please provide a valid attendance scan payload.", err.Error())
	}

	if len(payload.Embeddings) == 0 {
		return response.Error(c, fiber.StatusBadRequest, "No palm was detected. Please try placing your hand again.", "embeddings cannot be empty")
	}

	input := service.ProcessAttendanceInput{
		DeviceCode:     payload.DeviceCode,
		ModelVersion:   payload.ModelVersion,
		EmbeddingDim:   payload.EmbeddingDim,
		Embedding:      payload.Embeddings[0], // Extract the first embedding
		LivenessPassed: payload.LivenessPassed,
		QualityScore:   payload.QualityScore,
		ThermalMin:     payload.ThermalMin,
		ThermalMax:     payload.ThermalMax,
		ThermalAvg:     payload.ThermalAvg,
	}

	result, err := h.attendanceSvc.ProcessPalmAttendance(input)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "We couldn't identify your palm. Please try placing your hand again.", err.Error())
	}

	return response.Success(c, fiber.StatusOK, result.Message, fiber.Map{
		"action": result.Action,
		"user": fiber.Map{
			"id":        result.UserID,
			"full_name": result.FullName,
		},
	})
}

func (h *DeviceHandler) IdentifyPalm(c *fiber.Ctx) error {
	var payload struct {
		DeviceCode     string      `json:"device_code"`
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
		return response.Error(c, fiber.StatusBadRequest, "Please provide a valid identification payload.", err.Error())
	}

	if len(payload.Embeddings) == 0 {
		return response.Error(c, fiber.StatusBadRequest, "No palm was detected. Please try placing your hand again.", "embeddings cannot be empty")
	}

	input := service.ProcessAttendanceInput{
		DeviceCode:     payload.DeviceCode,
		ModelVersion:   payload.ModelVersion,
		EmbeddingDim:   payload.EmbeddingDim,
		Embedding:      payload.Embeddings[0],
		LivenessPassed: payload.LivenessPassed,
		QualityScore:   payload.QualityScore,
		ThermalMin:     payload.ThermalMin,
		ThermalMax:     payload.ThermalMax,
		ThermalAvg:     payload.ThermalAvg,
	}

	result, err := h.attendanceSvc.IdentifyPalm(input)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "We couldn't identify your palm. Please try placing your hand again.", err.Error())
	}

	return response.Success(c, fiber.StatusOK, result.Message, fiber.Map{
		"user": fiber.Map{
			"id":        result.UserID,
			"full_name": result.FullName,
		},
	})
}
