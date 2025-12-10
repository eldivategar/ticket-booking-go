package event

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
	return &repository{
		db: db,
	}
}

func (r *repository) CreateEvent(ctx context.Context, event domain.Event) (*domain.Event, error) {
	if err := r.db.WithContext(ctx).Create(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *repository) GetEventByID(ctx context.Context, eventID uuid.UUID) (*domain.Event, error) {
	var event domain.Event
	err := r.db.WithContext(ctx).Where("id = ?", eventID).First(&event).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (r *repository) GetEventByName(ctx context.Context, eventName string) (*domain.Event, error) {
	var event domain.Event
	err := r.db.WithContext(ctx).Where("name = ?", eventName).First(&event).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (r *repository) GetAllEvent(ctx context.Context) ([]domain.Event, error) {
	var events []domain.Event
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Find(&events).
		Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (r *repository) UpdateEvent(ctx context.Context, event domain.Event) (*domain.Event, error) {
	if err := r.db.WithContext(ctx).Save(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *repository) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Event{}, "id = ?", eventID).Error; err != nil {
		return err
	}
	return nil
}
