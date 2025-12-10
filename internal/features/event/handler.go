package event

import (
	"go-war-ticket-service/internal/domain"
	"go-war-ticket-service/internal/platform/responses"
	"go-war-ticket-service/internal/platform/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	usecase   Usecase
	validator *validator.Validator
}

func NewHandler(uc Usecase, validator *validator.Validator) *Handler {
	return &Handler{
		usecase:   uc,
		validator: validator,
	}
}

func (h *Handler) CreateEvent(c *fiber.Ctx) error {
	var req EventRequest
	if err := c.BodyParser(&req); err != nil {
		return responses.Error(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.validator.Validate(req); err != nil {
		errors := h.validator.FormatErrors(err)
		return responses.ValidationError(c, errors)
	}

	res, err := h.usecase.CreateEvent(c.Context(), domain.Event{
		Name:           req.Name,
		Description:    req.Description,
		Location:       req.Location,
		Price:          req.Price,
		TotalStock:     req.TotalStock,
		AvailableStock: req.TotalStock,
		Image:          req.Image,
		Date:           req.Date,
	})
	if err != nil {
		return responses.UsecaseError(c, err)
	}

	response := EventResponse{
		ID:             res.ID,
		Name:           res.Name,
		Description:    res.Description,
		Location:       res.Location,
		Price:          res.Price,
		TotalStock:     res.TotalStock,
		AvailableStock: res.AvailableStock,
		Image:          res.Image,
		Date:           res.Date,
		CreatedAt:      res.CreatedAt,
		UpdatedAt:      res.UpdatedAt,
	}

	return responses.Success(c, response, "success")
}

func (h *Handler) GetEventByID(c *fiber.Ctx) error {
	eventIDParams := c.Params("event_id")
	eventID, err := uuid.Parse(eventIDParams)
	if err != nil {
		return responses.Error(c, fiber.StatusBadRequest, domain.ErrInvalidID.Error())
	}

	res, err := h.usecase.GetEventByID(c.Context(), eventID)
	if err != nil {
		return responses.UsecaseError(c, err)
	}

	response := EventResponse{
		ID:             res.ID,
		Name:           res.Name,
		Description:    res.Description,
		Location:       res.Location,
		Price:          res.Price,
		TotalStock:     res.TotalStock,
		AvailableStock: res.AvailableStock,
		Image:          res.Image,
		Date:           res.Date,
		CreatedAt:      res.CreatedAt,
		UpdatedAt:      res.UpdatedAt,
	}

	return responses.Success(c, response, "success")
}
func (h *Handler) GetAllEvent(c *fiber.Ctx) error {
	res, err := h.usecase.GetAllEvent(c.Context())
	if err != nil {
		return responses.UsecaseError(c, err)
	}

	response := make([]EventResponse, len(res))
	for i, event := range res {
		response[i] = EventResponse{
			ID:             event.ID,
			Name:           event.Name,
			Description:    event.Description,
			Location:       event.Location,
			Price:          event.Price,
			TotalStock:     event.TotalStock,
			AvailableStock: event.AvailableStock,
			Image:          event.Image,
			Date:           event.Date,
			CreatedAt:      event.CreatedAt,
			UpdatedAt:      event.UpdatedAt,
		}
	}

	return responses.Success(c, response, "success")
}

func (h *Handler) DeleteEvent(c *fiber.Ctx) error {
	eventIDParams := c.Params("event_id")
	eventID, err := uuid.Parse(eventIDParams)
	if err != nil {
		return responses.Error(c, fiber.StatusBadRequest, domain.ErrInvalidID.Error())
	}

	if err := h.usecase.DeleteEvent(c.Context(), eventID); err != nil {
		return responses.UsecaseError(c, err)
	}

	return responses.Success(c, nil, "success")
}
