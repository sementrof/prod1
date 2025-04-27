package models

import "errors"

var (
	ErrConnectionDb           = errors.New("error connection to database")
	ErrRequestPayload         = errors.New("invalid request payload")
	ErrLoginOrPasswordOrEmail = errors.New("missing required fields email or password or user")
)
