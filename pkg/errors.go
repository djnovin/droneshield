package pkg

import (
	"encoding/json"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func (e *AppError) Send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	json.NewEncoder(w).Encode(e)
}

var (
  // ErrNotFound is returned when a resource is not found
  ErrNotFound = NewAppError(http.StatusNotFound, "Resource not found")
  // ErrBadRequest is returned when the request is invalid
  ErrBadRequest = NewAppError(http.StatusBadRequest, "Bad request")
  // ErrInternalServer is returned when an internal server error occurs
  ErrInternal = NewAppError(http.StatusInternalServerError, "Internal server error")
)
