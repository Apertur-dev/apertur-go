package apertur

import "fmt"

// AperturError is the base error type for all API errors returned by the Apertur service.
type AperturError struct {
	// StatusCode is the HTTP status code from the response.
	StatusCode int `json:"status_code"`
	// Code is the machine-readable error code from the API (e.g. "NOT_FOUND").
	Code string `json:"code"`
	// Message is the human-readable error description.
	Message string `json:"message"`
}

// Error implements the error interface.
func (e *AperturError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("apertur: HTTP %d [%s] %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("apertur: HTTP %d %s", e.StatusCode, e.Message)
}

// AuthenticationError indicates a 401 Unauthorized response.
type AuthenticationError struct {
	AperturError
}

// NewAuthenticationError creates an AuthenticationError with the given message.
func NewAuthenticationError(message string) *AuthenticationError {
	return &AuthenticationError{
		AperturError: AperturError{
			StatusCode: 401,
			Code:       "AUTHENTICATION_FAILED",
			Message:    message,
		},
	}
}

// NotFoundError indicates a 404 Not Found response.
type NotFoundError struct {
	AperturError
}

// NewNotFoundError creates a NotFoundError with the given message.
func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		AperturError: AperturError{
			StatusCode: 404,
			Code:       "NOT_FOUND",
			Message:    message,
		},
	}
}

// RateLimitError indicates a 429 Too Many Requests response.
type RateLimitError struct {
	AperturError
	// RetryAfter is the number of seconds to wait before retrying, if provided by the server.
	RetryAfter int `json:"retry_after"`
}

// NewRateLimitError creates a RateLimitError with the given message and retry-after duration.
func NewRateLimitError(message string, retryAfter int) *RateLimitError {
	return &RateLimitError{
		AperturError: AperturError{
			StatusCode: 429,
			Code:       "RATE_LIMIT",
			Message:    message,
		},
		RetryAfter: retryAfter,
	}
}

// ValidationError indicates a 400 Bad Request response due to invalid input.
type ValidationError struct {
	AperturError
}

// NewValidationError creates a ValidationError with the given message.
func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		AperturError: AperturError{
			StatusCode: 400,
			Code:       "VALIDATION_ERROR",
			Message:    message,
		},
	}
}
