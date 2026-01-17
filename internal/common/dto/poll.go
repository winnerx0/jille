package dto

import "time"

type CreatePollRequest struct {
	Title     string    `json:"title" validate:"required"`
	Options   []string  `json:"options" validate:"required,optionlistmin=2,optionlistmax=15"`
	ExpiresAt time.Time `json:"expires_at" validate:"required"`
}

type PollResponse struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Options []string `json:"options"`
}

type PollViewResponse struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Options []Option `json:"options"`
}

type Option struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"voter_count"`
}
