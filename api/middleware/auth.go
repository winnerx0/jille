package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/winnerx0/jille/internal/application"
)

func JWTMiddleware(c fiber.Ctx, jwtservice application.JwtService) error {

	authorization := c.Get("Authorization")
	if authorization == "" {

		c.Response().SetStatusCode(401)
		return c.JSON(fiber.Map{"message": "Authorization header required"})
	}

	parts := strings.SplitN(authorization, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {

		c.Response().SetStatusCode(401)
		return c.JSON(fiber.Map{"message": "Invalid token provided"})
	}

	token := parts[1]

	isVerified, err := jwtservice.VerifyAccessToken(token)

	if err != nil || !isVerified {
		c.Response().SetStatusCode(401)
		return c.JSON(fiber.Map{"message": "Invalid token provided"})
	}

	userID, err := jwtservice.GetTokenSubject(token)

	if err != nil {

		c.Response().SetStatusCode(401)
		return c.JSON(fiber.Map{"message": "Invalid token provided"})
	}

	c.Locals("userID", userID)

	return c.Next()
}
