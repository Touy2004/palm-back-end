package response

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Success formats a standard successful API response.
func Success(c *fiber.Ctx, code int, message string, data interface{}) error {
	return c.Status(code).JSON(fiber.Map{
		"code":    code,
		"status":  http.StatusText(code),
		"message": message,
		"data":    data,
	})
}

// SuccessWithMeta formats a standard successful API response with additional meta info (like pagination).
func SuccessWithMeta(c *fiber.Ctx, code int, message string, data interface{}, meta interface{}) error {
	return c.Status(code).JSON(fiber.Map{
		"code":    code,
		"status":  http.StatusText(code),
		"message": message,
		"data":    data,
		"meta":    meta,
	})
}

// Error formats a standard error API response.
func Error(c *fiber.Ctx, code int, message string, err interface{}) error {
	return c.Status(code).JSON(fiber.Map{
		"code":    code,
		"status":  http.StatusText(code),
		"message": message,
		"error":   err,
	})
}
