package order

import (
	"time"

	"github.com/google/uuid"
)

type OrderRequest struct {
	EventID  uuid.UUID `json:"event_id" validate:"required"`
	Quantity int       `json:"quantity" validate:"required,numeric,min=1"`
}

type OrderResponse struct {
	BookingID string    `json:"booking_id"`
	Event     Event     `json:"event"`
	Quantity  int       `json:"quantity"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Tickets   []string  `json:"tickets,omitempty"`
}

type Event struct {
	Name     string    `json:"name"`
	Location string    `json:"location"`
	Date     time.Time `json:"date"`
	Image    string    `json:"image"`
}

type PaymentWebhookRequest struct {
	BookingID     string `json:"booking_id" validate:"required"`
	PaymentStatus string `json:"payment_status" validate:"required"`
}
