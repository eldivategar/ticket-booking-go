package auth

import (
	"go-service-boilerplate/internal/platform/response"
	"go-service-boilerplate/internal/platform/validator"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	usecase  Usecase
	validate *validator.Validator
}

func NewHandler(uc Usecase, val *validator.Validator) *Handler {
	return &Handler{
		usecase:  uc,
		validate: val,
	}
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	if err := h.validate.Validate(&req); err != nil {
		errors := h.validate.FormatErrors(err)
		return response.ValidationError(c, errors)
	}

	res, err := h.usecase.Register(c.Context(), req)
	if err != nil {
		return response.UsecaseError(c, err)
	}

	return response.Success(c, res, "registration successful")
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	if err := h.validate.Validate(&req); err != nil {
		errors := h.validate.FormatErrors(err)
		return response.ValidationError(c, errors)
	}

	res, err := h.usecase.Login(c.Context(), req)
	if err != nil {
		return response.UsecaseError(c, err)
	}

	return response.Success(c, res, "login successful")
}
