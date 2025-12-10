package user

import (
	"go-war-ticket-service/internal/domain"
	"go-war-ticket-service/internal/platform/responses"
	"go-war-ticket-service/internal/utils/contextutil"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	usecase Usecase
}

func NewHandler(uc Usecase) *Handler {
	return &Handler{
		usecase: uc,
	}
}

func (h *Handler) GetMyProfile(c *fiber.Ctx) error {
	userId, err := contextutil.GetUserID(c.Context())
	if err != nil {
		return responses.Error(c, fiber.StatusInternalServerError, domain.ErrInternal.Error())
	}

	res, err := h.usecase.GetProfile(c.Context(), userId)
	if err != nil {
		return responses.UsecaseError(c, err)
	}

	response := UserResponse{
		ID:       res.ID,
		FullName: res.FullName,
		Username: res.Username,
		Email:    res.Email,
		Avatar:   res.Avatar,
	}

	return responses.Success(c, response, "success")
}
