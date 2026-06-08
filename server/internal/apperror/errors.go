package apperror

import "net/http"

type AppError struct {
	Code    string
	Message string
	Status  int
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func Unauthorized(message string) *AppError {
	return &AppError{Code: "UNAUTHORIZED", Message: message, Status: http.StatusUnauthorized}
}

func Forbidden(message string) *AppError {
	return &AppError{Code: "FORBIDDEN", Message: message, Status: http.StatusForbidden}
}

func BadRequest(message string) *AppError {
	return &AppError{Code: "INVALID_PAYLOAD", Message: message, Status: http.StatusBadRequest}
}

func Internal(message string, err error) *AppError {
	return &AppError{Code: "INTERNAL_ERROR", Message: message, Status: http.StatusInternalServerError, Err: err}
}
