package order

import (
	"go-war-ticket-service/internal/domain"
	"go-war-ticket-service/internal/platform/responses"
	"go-war-ticket-service/internal/platform/validator"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	usecase   Usecase
	service   Service
	validator *validator.Validator
}

func NewHandler(uc Usecase, svc Service, validator *validator.Validator) *Handler {
	return &Handler{
		usecase:   uc,
		service:   svc,
		validator: validator,
	}
}

func (h *Handler) CreateOrder(c *fiber.Ctx) error {
	var orderReq OrderRequest
	if err := c.BodyParser(&orderReq); err != nil {
		return responses.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	if err := h.validator.Validate(orderReq); err != nil {
		errors := h.validator.FormatErrors(err)
		return responses.ValidationError(c, errors)
	}

	order := domain.Order{
		EventID:  orderReq.EventID,
		Quantity: orderReq.Quantity,
	}

	createdOrder, err := h.usecase.CreateOrder(c.Context(), order)
	if err != nil {
		return responses.UsecaseError(c, err)
	}

	response := OrderResponse{
		BookingID: createdOrder.BookingID,
		Event: Event{
			Name:     createdOrder.Event.Name,
			Location: createdOrder.Event.Location,
			Date:     createdOrder.Event.Date,
			Image:    createdOrder.Event.Image,
		},
		Quantity:  createdOrder.Quantity,
		Total:     createdOrder.Event.Price * float64(createdOrder.Quantity),
		Status:    string(createdOrder.Status),
		CreatedAt: createdOrder.CreatedAt,
	}

	return responses.Success(c, response, "Order created successfully")
}

func (h *Handler) GetOrderByBookingID(c *fiber.Ctx) error {
	bookingID := c.Params("booking_id")
	if bookingID == "" {
		return responses.Error(c, fiber.StatusBadRequest, "Invalid booking ID")
	}

	order, err := h.usecase.GetOrderByBookingID(c.Context(), bookingID)
	if err != nil {
		return responses.UsecaseError(c, err)
	}

	response := OrderResponse{
		BookingID: order.BookingID,
		Event: Event{
			Name:     order.Event.Name,
			Location: order.Event.Location,
			Date:     order.Event.Date,
			Image:    order.Event.Image,
		},
		Quantity:  order.Quantity,
		Total:     order.Event.Price * float64(order.Quantity),
		Status:    string(order.Status),
		CreatedAt: order.CreatedAt,
		Tickets:   []string{},
	}

	for _, ticket := range order.Ticket {
		response.Tickets = append(response.Tickets, ticket.PDFUrl)
	}

	return responses.Success(c, response, "Order retrieved successfully")
}

func (h *Handler) GetOrderList(c *fiber.Ctx) error {
	orders, err := h.usecase.GetOrderList(c.Context())
	if err != nil {
		return responses.UsecaseError(c, err)
	}

	response := make([]OrderResponse, len(orders))
	for i, order := range orders {
		response[i] = OrderResponse{
			BookingID: order.BookingID,
			Event: Event{
				Name:     order.Event.Name,
				Location: order.Event.Location,
				Date:     order.Event.Date,
				Image:    order.Event.Image,
			},
			Quantity:  order.Quantity,
			Total:     order.Event.Price * float64(order.Quantity),
			Status:    string(order.Status),
			CreatedAt: order.CreatedAt,
		}
	}

	return responses.Success(c, response, "Orders retrieved successfully")
}

func (h *Handler) ProcessPaymentWebhook(c *fiber.Ctx) error {
	var payload PaymentWebhookRequest
	if err := c.BodyParser(&payload); err != nil {
		return responses.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	if err := h.service.ProcessPaymentWebhook(c.Context(), payload); err != nil {
		return responses.UsecaseError(c, err)
	}

	return responses.Success(c, nil, "Payment processed successfully")
}
