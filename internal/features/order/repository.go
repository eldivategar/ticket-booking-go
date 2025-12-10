package order

import (
	"context"
	"go-war-ticket-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateOrder(ctx context.Context, order *domain.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create order
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		result := tx.Model(&domain.Event{}).
			Where("id = ?", order.EventID).
			UpdateColumn("available_stock", gorm.Expr("available_stock - ?", order.Quantity))

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return domain.ErrEventNotFound
		}

		if err := tx.First(&order.Event, order.EventID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *repository) GetOrderByBookingID(ctx context.Context, bookingID string) (*domain.Order, error) {
	var order domain.Order

	if err := r.db.Model(&domain.Order{}).
		Where("booking_id = ?", bookingID).
		Preload("Event").
		Preload("Ticket").
		Find(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *repository) GetOrderList(ctx context.Context, userID uuid.UUID) ([]domain.Order, error) {
	var orders []domain.Order

	if err := r.db.Model(&domain.Order{}).
		Where("user_id = ?", userID).
		Preload("Event").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *repository) GetEventByID(ctx context.Context, eventID uuid.UUID) (*domain.Event, error) {
	var event domain.Event

	if err := r.db.Model(&domain.Event{}).Where("id = ?", eventID).Scan(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *repository) UpdateOrderStatus(ctx context.Context, bookingID string, status domain.OrderStatus) error {
	return r.db.Model(&domain.Order{}).
		Where("booking_id = ?", bookingID).
		Update("status", status).Error
}
