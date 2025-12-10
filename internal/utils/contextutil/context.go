package contextutil

import (
	"context"
	"errors"
	"go-war-ticket-service/internal/utils"

	"github.com/google/uuid"
)

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(utils.UserID).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user ID not found in context")
	}

	return id, nil
}
