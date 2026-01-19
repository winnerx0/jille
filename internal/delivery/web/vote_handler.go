package web

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/winnerx0/jille/internal/application"
	"github.com/winnerx0/jille/internal/common/dto"
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

func (h *votehandler) VotePoll(b utils.Broker) fiber.Handler {

	return func(c fiber.Ctx) error {
		var voteRequst dto.VoteRequest

		if err := c.Bind().Body(&voteRequst); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "Failed to parse body"})
		}

		ctx := context.WithValue(c.Context(), "userID", c.Locals("userID"))
		response, err := h.voteservice.VotePoll(ctx, voteRequst)

		if err != nil {
			if errors.Is(err, utils.PollExpiredError) {
				return c.Status(400).JSON(fiber.Map{"message": err.Error()})
			} else if errors.Is(err, utils.OptionNotFound) {
				return c.Status(404).JSON(fiber.Map{"message": err.Error()})
			} else {
				return c.Status(500).JSON(fiber.Map{"message": err.Error()})
			}
		}

		b.Events <- utils.Event{
			Type:    "POLL_VOTE",
			Payload: voteRequst,
		}

		return c.JSON(response)
	}
}
