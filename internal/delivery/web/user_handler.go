package web

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/application"
)

type UserHandler struct {
	userservice application.UserService
}

func NewUserHandler(userservice application.UserService) *UserHandler {
	return &UserHandler{
		userservice: userservice,
	}
}

func (h *UserHandler) GetUser(c fiber.Ctx) error {

	userID := c.Params("userID")

	response, err := h.userservice.GetUserById(c.RequestCtx(), uuid.MustParse(userID))

	if err != nil {
		c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User retrieved successfully", "data": response})
}

// fiber:context-methods migrated
