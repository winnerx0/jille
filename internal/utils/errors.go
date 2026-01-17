package utils

import "errors"

var (
	UserExistsError    = errors.New("User with email already exists")
	UserNotFoundError  = errors.New("User not found")
	TokenExpiredError  = errors.New("Token expired")
	TokenNotFoundError = errors.New("Token not found")
	PollExpiredError   = errors.New("Poll has expired")
	OptionNotFound     = errors.New("Option not found in poll")
	PollNotFoundError  = errors.New("Poll not found")
	PollAccessDeniedError = errors.New("Only the creator can view the live votings")
)
