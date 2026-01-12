package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	userservice UserService
}

func NewUserHandler(userservice UserService) *UserHandler {
	return &UserHandler{
		userservice: userservice,
	}
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {

	userID := c.Params("userID")

	response, err := h.userservice.GetUserById(c.Context(), uuid.MustParse(userID))

	if err != nil {
		c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User retrieved successfully", "data": response})
}
