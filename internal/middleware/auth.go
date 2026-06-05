package middleware

import (
	"strings"

	"github.com/Touy2004/palm-back-end/pkg/jwt"
	"github.com/Touy2004/palm-back-end/pkg/response"

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
		return response.Error(c, fiber.StatusUnauthorized, "Missing authorization header", "unauthorized")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == authHeader {
		return response.Error(c, fiber.StatusUnauthorized, "Invalid token format", "invalid token format")
	}

	claims, err := m.jwt.Parse(tokenStr)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "Invalid or expired token", err.Error())
	}

	c.Locals("user", claims)
	return c.Next()
}
