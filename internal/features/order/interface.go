package order

import (
	"context"
	"go-war-ticket-service/internal/domain"

	"github.com/google/uuid"
)

type Usecase interface {
	CreateOrder(ctx context.Context, order domain.Order) (*domain.Order, error)
	GetOrderByBookingID(ctx context.Context, bookingID string) (*domain.Order, error)
	GetOrderList(ctx context.Context) ([]domain.Order, error)
}

type Repository interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	GetEventByID(ctx context.Context, eventID uuid.UUID) (*domain.Event, error)
	GetOrderByBookingID(ctx context.Context, bookingID string) (*domain.Order, error)
	GetOrderList(ctx context.Context, userID uuid.UUID) ([]domain.Order, error)
	UpdateOrderStatus(ctx context.Context, bookingID string, status domain.OrderStatus) error
}

type Service interface {
	ProcessPaymentWebhook(ctx context.Context, payload PaymentWebhookRequest) error
}
