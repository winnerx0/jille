package auth

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/winnerx0/jille/internal/common/dto"
)

type authHandler struct {
	authservie Service
}

func NewAuthHandler(authservice Service) *authHandler {
	return &authHandler{
		authservie: authservice,
	}
}

func (h *authHandler) RegisterUser(c *fiber.Ctx) error {

	var registerRequest dto.CreateUserRequest

	err := json.Unmarshal(c.Body(), &registerRequest)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid credentials format"})
	}

	response, err := h.authservie.Register(c.Context(), registerRequest)

	if err != nil {
		fmt.Println("error", err.Error())
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	// if response != nil {
	// 	return c.Status(404).JSON(fiber.Map{"message": "User not found"})
	// }

	return c.JSON(fiber.Map{"message": "Registration Successful", "data": response})

}

func (h *authHandler) LoginUser(c *fiber.Ctx) error {

	var loginRequest dto.LoginUserRequest

	err := json.Unmarshal(c.Body(), &loginRequest)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid credentials format"})
	}

	response, err := h.authservie.Login(c.Context(), loginRequest)

	if err != nil {
		fmt.Println("error", err.Error())
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	if response == nil {
		return c.Status(404).JSON(fiber.Map{"message": "User not found"})
	}

	return c.JSON(fiber.Map{"message": "Login Successful", "data": response})

}
