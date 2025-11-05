package response

import "github.com/gofiber/fiber/v2"

// APIResponse represents a standard structure for API responses
type APIResponse struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

func Success(c *fiber.Ctx, data any, message string) error {
	if message == "" {
		// Default success message
		message = "success"
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Success:    true,
		StatusCode: fiber.StatusOK,
		Message:    message,
		Data:       data,
	})
}

func Error(c *fiber.Ctx, statusCode int, message string) error {
	if message == "" {
		// Default error message
		message = "an error occurred"
	}

	return c.Status(statusCode).JSON(APIResponse{
		Success:    false,
		StatusCode: statusCode,
		Message:    message,
		Data:       nil,
	})
}
