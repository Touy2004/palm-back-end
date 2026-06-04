package middleware

import (
	"github.com/Touy2004/palm-back-end/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

type RoleMiddleware struct{}

func NewRoleMiddleware() *RoleMiddleware {
	return &RoleMiddleware{}
}

func (m *RoleMiddleware) Require(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("user").(*jwt.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		hasRole := false
		for _, role := range roles {
			if claims.Role == role {
				hasRole = true
				break
			}
		}

		if hasRole {
			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
}
