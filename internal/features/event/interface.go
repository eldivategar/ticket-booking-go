package event

import (
	"context"
	"go-war-ticket-service/internal/domain"

	"github.com/google/uuid"
)

type Usecase interface {
	CreateEvent(ctx context.Context, event domain.Event) (*domain.Event, error)
	GetEventByID(ctx context.Context, eventID uuid.UUID) (*domain.Event, error)
	GetAllEvent(ctx context.Context) ([]domain.Event, error)
	DeleteEvent(ctx context.Context, eventID uuid.UUID) error
}

type Repository interface {
	CreateEvent(ctx context.Context, event domain.Event) (*domain.Event, error)
	GetEventByID(ctx context.Context, eventID uuid.UUID) (*domain.Event, error)
	GetEventByName(ctx context.Context, eventName string) (*domain.Event, error)
	GetAllEvent(ctx context.Context) ([]domain.Event, error)
	UpdateEvent(ctx context.Context, event domain.Event) (*domain.Event, error)
	DeleteEvent(ctx context.Context, eventID uuid.UUID) error
}
