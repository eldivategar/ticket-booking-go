package ticket

import (
	"context"
	"go-war-ticket-service/internal/domain"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateTicket(ctx context.Context, ticket *domain.Ticket) error {
	return r.db.Create(&ticket).Error
}