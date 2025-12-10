package domain

import "github.com/google/uuid"

type TicketStatus string

const (
	TicketStatusValid TicketStatus = "VALID"
	TicketStatusUsed  TicketStatus = "USED"
)

type Ticket struct {
	BaseModel
	OrderID uuid.UUID `gorm:"not null" json:"order_id"`
	EventID uuid.UUID `gorm:"not null" json:"event_id"`
	UserID  uuid.UUID `gorm:"not null" json:"user_id"`

	TicketNumber string       `gorm:"type:varchar(50);not null;uniqueIndex" json:"ticket_number"`
	PDFUrl       string       `gorm:"type:text" json:"pdf_url"`
	// Status       TicketStatus `gorm:"type:varchar(50);default:'VALID';index" json:"status"` // VALID, USED
}
