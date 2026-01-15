package application

import "github.com/winnerx0/jille/internal/application/repository"

type voteservice struct {

	repo repository.VoteRepository
}

func NewVoteService(repo repository.VoteRepository)