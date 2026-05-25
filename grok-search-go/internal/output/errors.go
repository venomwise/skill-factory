package output

import (
	"errors"

	"github.com/venomwise/skill-factory/grok-search/internal/client"
	"github.com/venomwise/skill-factory/grok-search/internal/config"
)

// CommandError is returned after a structured runtime error has been rendered.
type CommandError struct {
	Code   string
	Detail string
	Err    error
}

func (e *CommandError) Error() string {
	if e.Detail != "" {
		return e.Code + ": " + e.Detail
	}
	return e.Code
}

func (e *CommandError) Unwrap() error {
	return e.Err
}

// NewCommandError creates a non-zero command error for an already-rendered response.
func NewCommandError(code, detail string, err error) *CommandError {
	return &CommandError{Code: code, Detail: detail, Err: err}
}

// NewErrorResponse creates a normalized runtime error response.
func NewErrorResponse(code, detail string) ErrorResponse {
	return ErrorResponse{OK: false, Error: code, Detail: detail}
}

// FromError maps known internal errors to normalized runtime error responses.
func FromError(err error) ErrorResponse {
	var cfgErr *config.Error
	if errors.As(err, &cfgErr) {
		return NewErrorResponse(cfgErr.Code, cfgErr.Detail)
	}

	var reqErr *client.RequestError
	if errors.As(err, &reqErr) {
		return NewErrorResponse("request_failed", reqErr.Detail)
	}

	var cmdErr *CommandError
	if errors.As(err, &cmdErr) {
		return NewErrorResponse(cmdErr.Code, cmdErr.Detail)
	}

	if err == nil {
		return NewErrorResponse("unknown_error", "")
	}
	return NewErrorResponse("request_failed", err.Error())
}
