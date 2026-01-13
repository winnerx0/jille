package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/application"
	"github.com/winnerx0/jille/internal/common/dto"
	"github.com/winnerx0/jille/internal/utils"
)

type pollhandler struct {
	pollservice application.PollService
	validator   utils.XValidator
}

func NewPollHandler(pollservice application.PollService, validator utils.XValidator) *pollhandler {
	return &pollhandler{
		pollservice: pollservice,
		validator:   validator,
	}
}

func (h *pollhandler) CreatePoll(c *fiber.Ctx) error {

	var pollRequest dto.CreatePollRequest


	if err := c.BodyParser(&pollRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

		if err := h.validator.Validate(pollRequest); err != nil {
		return c.Status(422).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	if err := h.pollservice.CreatePoll(c.Context(), &pollRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Poll created successfully"})
}

func (h *pollhandler) DeletePoll(c *fiber.Ctx) error {

	pollID := c.Params("pollID")

	if err := h.pollservice.DeletePoll(c.Context(), uuid.MustParse(pollID)); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Poll deleted successfully"})
}