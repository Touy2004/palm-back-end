package middleware

import (
	"fmt"
	"slices"

	"github.com/Touy2004/palm-back-end/pkg/jwt"
	"github.com/Touy2004/palm-back-end/pkg/response"

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
			return response.Error(c, fiber.StatusUnauthorized, "Unauthorized access", "unauthorized")
		}

		if slices.Contains(roles, claims.Role) {
			fmt.Println("Role: ", claims.Role)
			return c.Next()
		}

		fmt.Println("Forbidden: ", roles)
		return response.Error(c, fiber.StatusForbidden, "Forbidden access", "forbidden")
	}
}
