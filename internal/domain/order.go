package domain

import "github.com/google/uuid"

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusPaid       OrderStatus = "PAID"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusCompleted  OrderStatus = "COMPLETED"
	OrderStatusFailed     OrderStatus = "FAILED"
)

type Order struct {
	BaseModel
	BookingID  string      `gorm:"type:varchar(25);not null;uniqueIndex" json:"booking_id"`
	UserID     uuid.UUID   `gorm:"not null" json:"user_id"`
	EventID    uuid.UUID   `gorm:"not null" json:"event_id"`
	Quantity   int         `gorm:"not null" json:"quantity"`
	// TotalPrice float64     `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Status     OrderStatus `gorm:"type:varchar(50);default:'PENDING';index" json:"status"`

	User   User     `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Event  Event    `gorm:"foreignKey:EventID;references:ID" json:"event"`
	Ticket []Ticket `gorm:"foreignKey:OrderID;references:ID" json:"ticket"`
}
