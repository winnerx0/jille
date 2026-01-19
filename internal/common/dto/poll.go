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
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Options   []Option  `json:"options"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatorID string    `json:"creator_id"`
	Voted     bool      `json:"voted"`
}

type Vote struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	PollID   string `json:"poll_id"`
	OptionID string `json:"option_id"`
}

type Option struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Votes []Vote `json:"votes"`
}

type ApiResponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}
