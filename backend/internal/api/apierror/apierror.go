package apierror

import (
	"fmt"
	"strings"
)

// APIError is a custom error type
// It implements the standard 'error' interface.
type APIError struct {
	StatusCode int
	Message    string
}

// Error makes APIError conform to the error interface.
func (e *APIError) Error() string {
	return e.Message
}

// New creates a new APIError.
func New(statusCode int, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// --- Helper functions for common error types ---

func NewNotFoundError(resource string) *APIError {
	return &APIError{
		StatusCode: 404,
		Message:    fmt.Sprintf("%s not found", resource),
	}
}

func NewForbiddenError() *APIError {
	return &APIError{
		StatusCode: 403,
		Message:    "permission denied",
	}
}

func NewBadRequestError(message string) *APIError {
	return &APIError{
		StatusCode: 400,
		Message:    message,
	}
}

func NewUnauthorizedError() *APIError {
	return &APIError{
		StatusCode: 401,
		Message:    "unauthorized",
	}
}

func NewInternalServerError(messages ...string) *APIError {
	message := "An unexpected internal server has occured:"

	if len(messages) > 0 {
		message = strings.Join(messages, " ")
	}

	return &APIError{
		StatusCode: 500,
		Message:    message,
	}
}
