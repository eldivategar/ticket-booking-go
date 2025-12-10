package ticket

import (
	"context"
	"go-war-ticket-service/internal/domain"
)

type Repository interface {
	CreateTicket(ctx context.Context, ticket *domain.Ticket) error
}
