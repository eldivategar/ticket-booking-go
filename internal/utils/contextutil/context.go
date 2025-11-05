package contextutil

import (
	"context"
	"errors"
	"go-service-boilerplate/internal/utils"

	"github.com/google/uuid"
)

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(utils.UserID).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user ID not found in context")
	}

	return id, nil
}
