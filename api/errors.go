package api

import (
	"net/http"
)

var (
	// ErrNotFound represents generic not found error
	ErrNotFound = SejastipError{
		Message:    "Tidak ditemukan",
		ErrorCode:  404,
		HTTPStatus: http.StatusNotFound,
	}

	// ErrInvalidParameter represents generic invalid param error
	ErrInvalidParameter = SejastipError{
		Message:    "Parameter tidak valid",
		ErrorCode:  401,
		HTTPStatus: http.StatusBadRequest,
	}
)

// SejastipError defines our custom error
type SejastipError struct {
	Message    string `json:"message"`
	ErrorCode  int    `json:"error_code"`
	HTTPStatus int    `json:"-"`
}

// Error returns the string representation of our custom error
// Satisfies the error interface
func (e SejastipError) Error() string {
	return e.Message
}
