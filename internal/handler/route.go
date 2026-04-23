package handler

import (
	"github.com/Touy2004/palm-back-end/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Auth *AuthHandler
}

type Middleware struct {
	Auth *middleware.AuthMiddleware
	Role *middleware.RoleMiddleware
}

func SetupRoutes(app *fiber.App, h *Handler, m *Middleware) {
	api := app.Group("/api")

	// Public routes
	auth := api.Group("/auth")
	auth.Post("/register", h.Auth.Register)
	auth.Post("/login", h.Auth.Login)

	// Authenticated routes
	user := api.Group("/user", m.Auth.Authenticate)
	user.Get("/profile", h.Auth.GetProfile)

	// Admin only routes
	admin := api.Group("/admin", m.Auth.Authenticate, m.Role.Require("user"))
	admin.Get("/users", h.Auth.GetUsers)
}
