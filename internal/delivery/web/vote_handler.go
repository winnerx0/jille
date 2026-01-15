package web

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/application"
	"github.com/winnerx0/jille/internal/utils"
)

type votehandler struct {
	voteservice application.VoteService
}

func NewVoteHandler(voteservice application.VoteService) *votehandler {
	return &votehandler{
		voteservice: voteservice,
	}
}

func (h *votehandler) VotePoll(c fiber.Ctx) error {

	pollID := uuid.MustParse(c.Params("pollId"))

	optionId := uuid.MustParse(c.Params("optionId"))

	response, err := h.voteservice.VotePoll(c.RequestCtx(), pollID, optionId)

	if err != nil {
		if errors.Is(err, utils.PollExpiredError) {
			return c.Status(400).JSON(fiber.Map{"message": err.Error()})
		} else if errors.Is(err, utils.OptionNotFound) {
			return c.Status(404).JSON(fiber.Map{"message": err.Error()})
		} else {
			return c.Status(500).JSON(fiber.Map{"message": err.Error()})
		}
	}

	return c.JSON(response)
}

// fiber:context-methods migrated
