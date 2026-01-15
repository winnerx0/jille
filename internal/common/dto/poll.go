package dto

import "time"

type CreatePollRequest struct {
	Title   string   `json:"title" validate:"required"`
	Options []string `json:"options" validate:"required,optionlistmin=2,optionlistmax=15"`
	ExpiresAt time.Time
}

type PollResponse struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Options []string `json:"options"`
}
