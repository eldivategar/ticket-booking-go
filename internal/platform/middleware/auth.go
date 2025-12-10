package middleware

import (
	"go-war-ticket-service/internal/domain"
	"go-war-ticket-service/internal/platform/responses"
	"go-war-ticket-service/internal/utils"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AuthRequired is a middleware that checks for a valid JWT token in the Authorization header
func AuthRequired(secret string, log *zap.SugaredLogger) fiber.Handler {
	return jwtware.New(jwtware.Config{
		// Use jwtware to validate the token
		SigningKey:  jwtware.SigningKey{Key: []byte(secret)},
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
		ContextKey:  "user",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Warnf("unauthorized access: %v", err)
			return responses.Error(c, fiber.StatusUnauthorized, domain.ErrUnauthorized.Error())
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			// Extract the token
			tokenData := c.Locals("user")
			if tokenData == nil {
				return responses.Error(c, fiber.StatusUnauthorized, domain.ErrUnauthorized.Error())
			}

			// Casting token
			token, ok := tokenData.(*jwt.Token)
			if !ok || !token.Valid {
				return responses.Error(c, fiber.StatusUnauthorized, domain.ErrUnauthorized.Error())
			}

			// Claims token
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return responses.Error(c, fiber.StatusUnauthorized, domain.ErrUnauthorized.Error())
			}

			// Get user ID in 'sub' claim
			userIDStr, ok := claims["sub"].(string)
			if !ok || userIDStr == "" {
				return responses.Error(c, fiber.StatusUnauthorized, domain.ErrUnauthorized.Error())
			}

			// Parse to UUID
			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				return responses.Error(c, fiber.StatusUnauthorized, domain.ErrUnauthorized.Error())
			}

			// Set user ID to context
			c.Locals(utils.UserID, userID)

			return c.Next()
		},
	})
}
