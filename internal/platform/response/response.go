package response

import (
	"go-service-boilerplate/internal/domain"

	"github.com/gofiber/fiber/v2"
)

// APIResponse represents a standard structure for API responses
type APIResponse struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

// Success sends a successful JSON response
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

// Error sends an error JSON response with the given status code and message
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

// ValidationError sends a validation error response with details
func ValidationError(c *fiber.Ctx, errors any) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(APIResponse{
		Success:    false,
		StatusCode: fiber.StatusUnprocessableEntity,
		Message:    "validation error",
		Data:       errors,
	})
}

// UsecaseError maps domain errors to appropriate HTTP responses
func UsecaseError(c *fiber.Ctx, err error) error {
	switch err {
	case domain.ErrEmailAlreadyExists:
		return Error(c, fiber.StatusBadRequest, err.Error())
	case domain.ErrUsernameAlreadyExists:
		return Error(c, fiber.StatusBadRequest, err.Error())
	case domain.ErrInvalidCredentials:
		return Error(c, fiber.StatusUnauthorized, err.Error())
	default:
		return Error(c, fiber.StatusInternalServerError, domain.ErrInternal.Error())
	}
}
