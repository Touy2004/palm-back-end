package handler

import (
	"github.com/Touy2004/palm-back-end/internal/constant"
	"github.com/Touy2004/palm-back-end/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Auth   *AuthHandler
	Admin  *AdminHandler
	User   *UserHandler
	Device *DeviceHandler
}

type Middleware struct {
	Auth *middleware.AuthMiddleware
	Role *middleware.RoleMiddleware
}

func SetupRoutes(app *fiber.App, h *Handler, m *Middleware) {
	api := app.Group("/api/v1")

	// Public routes
	auth := api.Group("/auth")
	auth.Post("/register", h.Auth.Register)
	auth.Post("/login", h.Auth.Login)
	auth.Post("/refresh", h.Auth.Refresh)

	// Authenticated routes
	user := api.Group("/me", m.Auth.Authenticate)
	user.Get("/", h.Auth.GetProfile)
	user.Get("/palm-templates", h.User.GetPalmTemplates)
	user.Delete("/palm-templates/:id", h.User.DeletePalmTemplate)
	user.Get("/attendance", h.User.GetMyAttendance)
	user.Patch("/password", h.User.ChangePassword)

	// Pairing routes (Requires authentication)
	pairing := api.Group("/pairing", m.Auth.Authenticate)
	pairing.Post("/scan", h.User.ScanPairingQR)
	pairing.Post("/approve", h.User.ApprovePairingQR)

	// Admin only routes
	admin := api.Group("/admin", m.Auth.Authenticate, m.Role.Require(constant.ROLE_ADMIN))
	
	// Admin Users (Search must be before :id to prevent collision in some routers)
	admin.Get("/users/search", h.Admin.SearchUsers)
	admin.Post("/users", h.Admin.CreateUser)
	admin.Get("/users", h.Admin.GetUsers)
	admin.Get("/users/:id", h.Admin.GetUserByID)
	admin.Patch("/users/:id", h.Admin.UpdateUser)
	admin.Delete("/users/:id", h.Admin.DeleteUser)

	// Admin Devices
	admin.Get("/devices", h.Admin.GetDevices)
	admin.Post("/devices", h.Admin.CreateDevice)
	admin.Patch("/devices/:id", h.Admin.UpdateDevice)

	// Admin Attendance
	admin.Get("/attendance", h.Admin.GetAttendanceHistory)
	admin.Get("/attendance/users/:user_id/history", h.Admin.GetUserAttendanceHistory)

	// Device endpoints (Hardware APIs)
	devices := api.Group("/devices")
	devices.Post("/heartbeat", h.Device.Heartbeat)
	devices.Post("/pairing-sessions", h.Device.CreatePairingSession)
	devices.Get("/pairing-sessions/:session_id/status", h.Device.GetSessionStatus)
	devices.Post("/palm/enroll", h.Device.EnrollPalm)
	devices.Post("/palm/identify", h.Device.IdentifyPalm)
	devices.Post("/attendance/palm", h.Device.ProcessAttendance)
}
