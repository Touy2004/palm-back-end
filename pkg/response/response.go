package response

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Success formats a standard successful API response.
func Success(c *fiber.Ctx, code int, message string, data any) error {
	return c.Status(code).JSON(fiber.Map{
		"code":    code,
		"status":  http.StatusText(code),
		"message": message,
		"data":    data,
	})
}

// SuccessWithMeta formats a standard successful API response with additional meta info (like pagination).
func SuccessWithMeta(c *fiber.Ctx, code int, message string, data any, meta any) error {
	return c.Status(code).JSON(fiber.Map{
		"code":    code,
		"status":  http.StatusText(code),
		"message": message,
		"data":    data,
		"meta":    meta,
	})
}

// Error formats a standard error API response.
func Error(c *fiber.Ctx, code int, message string, err any) error {
	if err != nil {
		// Log the actual technical error to the backend terminal for debugging
		log.Printf("[API ERROR] %s: %v\n", message, err)
	}

	return c.Status(code).JSON(fiber.Map{
		"code":    code,
		"status":  http.StatusText(code),
		"message": message,
	})
}
