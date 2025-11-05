package middleware

import (
	"go-service-boilerplate/internal/domain"
	"go-service-boilerplate/internal/platform/response"
	"go-service-boilerplate/internal/utils"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthRequired membuat middleware JWT baru
func AuthRequired(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check for Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, domain.ErrUnauthorized.Error())
		}

		// Use jwtware to validate the token
		jwtware.New(jwtware.Config{
			SigningKey:  jwtware.SigningKey{Key: []byte(secret)},
			TokenLookup: "header:Authorization",
			AuthScheme:  "Bearer",
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				return response.Error(c, fiber.StatusUnauthorized, domain.ErrUnauthorized.Error())
			},
			ContextKey: string(utils.UserID),
		})

		// Extract the token
		tokenData := c.Locals("user")
		if tokenData == nil {
			return response.Error(c, fiber.StatusUnauthorized, domain.ErrUnauthorized.Error())
		}

		// Casting token
		token, ok := tokenData.(*jwt.Token)
		if !ok || !token.Valid {
			return response.Error(c, fiber.StatusUnauthorized, domain.ErrUnauthorized.Error())
		}

		// Claims token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return response.Error(c, fiber.StatusUnauthorized, domain.ErrUnauthorized.Error())
		}

		// Get user ID in 'sub' claim
		userIDStr, ok := claims["sub"].(string)
		if !ok || userIDStr == "" {
			return response.Error(c, fiber.StatusUnauthorized, domain.ErrUnauthorized.Error())
		}

		// Parse to UUID
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, domain.ErrUnauthorized.Error())
		}

		// Set user ID to context
		c.Locals(string(utils.UserID), userID)

		return c.Next()
	}
}
