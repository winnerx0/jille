package dto

type VoteResponse struct {
	Message string
}

type VoteRequest struct {
	PollID string `json:"poll_id" validate:"required"`

	OptionID string `json:"option_id" validate:"required"`
}
