package web

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/winnerx0/jille/internal/application"
	"github.com/winnerx0/jille/internal/common/dto"
	"github.com/winnerx0/jille/internal/utils"
)

type authHandler struct {
	authservie application.AuthService
	validator  utils.XValidator
}

func NewAuthHandler(authservice application.AuthService, validator utils.XValidator) *authHandler {
	return &authHandler{
		authservie: authservice,
		validator:  validator,
	}
}

func (h *authHandler) RegisterUser(c fiber.Ctx) error {

	var registerRequest dto.CreateUserRequest

	if err := c.Bind().Body(&registerRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	if err := h.validator.Validate(registerRequest); err != nil {
		return c.Status(422).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	response, err := h.authservie.Register(c.RequestCtx(), registerRequest)

	if err != nil {
		fmt.Println("error", err.Error())
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	// if response != nil {
	// 	return c.Status(404).JSON(fiber.Map{"message": "User not found"})
	// }

	return c.JSON(response)

}

func (h *authHandler) LoginUser(c fiber.Ctx) error {

	var loginRequest dto.LoginUserRequest

	err := json.Unmarshal(c.Body(), &loginRequest)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid credentials format"})
	}

	if err := h.validator.Validate(loginRequest); err != nil {
		return c.Status(422).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	response, err := h.authservie.Login(c.RequestCtx(), loginRequest)

	if err != nil {
		fmt.Println("error", err.Error())
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	if response == nil {
		return c.Status(404).JSON(fiber.Map{"message": "User not found"})
	}

	return c.JSON(response)

}

func (h *authHandler) RefreshToken(c fiber.Ctx) error {

	var refreshTokenRequest dto.RefreshTokenRequest

	if err := c.Bind().Body(&refreshTokenRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	if err := h.validator.Validate(refreshTokenRequest); err != nil {
		return c.Status(422).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	response, err := h.authservie.RefreshToken(c.RequestCtx(), refreshTokenRequest)

	if err != nil {
		if err == utils.TokenNotFoundError {
			return c.Status(404).JSON(fiber.Map{"message": err.Error()})
		}
		if err == utils.TokenExpiredError {
			return c.Status(401).JSON(fiber.Map{"message": err.Error()})
		}
		fmt.Println("error", err.Error())
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	if response == nil {
		return c.Status(404).JSON(fiber.Map{"message": "User not found"})
	}

	return c.JSON(response)

}