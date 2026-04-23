package middleware

import (
	"strings"

	"github.com/Touy2004/palm-back-end/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	jwt *jwt.JWT
}

func NewAuthMiddleware(jwt *jwt.JWT) *AuthMiddleware {
	return &AuthMiddleware{jwt: jwt}
}

func (m *AuthMiddleware) Authenticate(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token format",
		})
	}

	claims, err := m.jwt.Parse(tokenStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid or expired token",
		})
	}

	c.Locals("user", claims)
	return c.Next()
}
