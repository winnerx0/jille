package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/winnerx0/jille/internal/domain/jwt"
)

func JWTMiddleware(c *fiber.Ctx, jwtservice jwt.JwtService) error {

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
	fmt.Println("parts", isVerified)

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
