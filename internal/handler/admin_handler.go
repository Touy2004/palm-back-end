package handler

import (
	"strconv"

	"github.com/Touy2004/palm-back-end/internal/model"
	"github.com/Touy2004/palm-back-end/internal/service"
	"github.com/Touy2004/palm-back-end/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// User Endpoints
func (h *AdminHandler) CreateUser(c *fiber.Ctx) error {
	var input service.RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	user, err := h.adminService.CreateUser(input)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Failed to create user", err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "User created successfully", user)
}

func (h *AdminHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.adminService.GetUsers()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch users", err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Users retrieved successfully", users)
}

func (h *AdminHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.adminService.GetUserByID(id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "User not found", err.Error())
	}
	return response.Success(c, fiber.StatusOK, "User retrieved successfully", user)
}

func (h *AdminHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var input map[string]interface{}
	if err := c.BodyParser(&input); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	user, err := h.adminService.UpdateUser(id, input)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update user", err.Error())
	}
	return response.Success(c, fiber.StatusOK, "User updated successfully", user)
}

func (h *AdminHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.adminService.DeleteUser(id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete user", err.Error())
	}
	return response.Success(c, fiber.StatusOK, "User deleted successfully", nil)
}

func (h *AdminHandler) SearchUsers(c *fiber.Ctx) error {
	q := c.Query("q")
	users, err := h.adminService.SearchUsers(q)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to search users", err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Users retrieved successfully", users)
}

// Device Endpoints
func (h *AdminHandler) GetDevices(c *fiber.Ctx) error {
	devices, err := h.adminService.GetDevices()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch devices", err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Devices retrieved successfully", devices)
}

func (h *AdminHandler) CreateDevice(c *fiber.Ctx) error {
	var device model.Device
	if err := c.BodyParser(&device); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	err := h.adminService.CreateDevice(&device)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create device", err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Device created successfully", device)
}

func (h *AdminHandler) UpdateDevice(c *fiber.Ctx) error {
	id := c.Params("id")
	var input map[string]interface{}
	if err := c.BodyParser(&input); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	device, err := h.adminService.UpdateDevice(id, input)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update device", err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Device updated successfully", device)
}

// Attendance Endpoints
func (h *AdminHandler) GetAttendanceHistory(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	logs, total, err := h.adminService.GetAttendanceHistory(page, limit, startDate, endDate)
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

func (h *AdminHandler) GetUserAttendanceHistory(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	logs, total, err := h.adminService.GetUserAttendanceHistory(userID, page, limit, startDate, endDate)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch user attendance history", err.Error())
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "User attendance history retrieved successfully", logs, fiber.Map{
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *AdminHandler) GetUserPalmTemplates(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	templates, err := h.adminService.GetUserPalmTemplates(userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch user palm templates", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "User palm templates retrieved successfully", templates)
}

func (h *AdminHandler) DeleteUserPalmTemplate(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	templateID := c.Params("template_id")

	err := h.adminService.DeleteUserPalmTemplate(userID, templateID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete user palm template", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "User palm template deleted successfully", nil)
}

func (h *AdminHandler) GetDashboardSummary(c *fiber.Ctx) error {
	summary, err := h.adminService.GetDashboardSummary()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get dashboard summary", err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Dashboard summary retrieved successfully", summary)
}