package event

import (
	"time"

	"github.com/google/uuid"
)

type EventRequest struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Location    string    `json:"location" validate:"required"`
	Price       float64   `json:"price" validate:"required"`
	TotalStock  int       `json:"total_stock" validate:"required"`
	Image       string    `json:"image" validate:"required"`
	Date        time.Time `json:"date" validate:"required"`
}

type EventResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Location       string    `json:"location"`
	Price          float64   `json:"price"`
	TotalStock     int       `json:"total_stock"`
	AvailableStock int       `json:"available_stock"`
	Image          string    `json:"image"`
	Date           time.Time `json:"date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
