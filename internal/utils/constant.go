package utils

import "github.com/google/uuid"

// For strong typing of user ID in context
type userID uuid.UUID

// UserID is the key used to store and retrieve user ID from context
var UserID userID
