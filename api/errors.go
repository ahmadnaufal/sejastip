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
		ErrorCode:  400,
		HTTPStatus: http.StatusBadRequest,
	}

	// ErrInvalidCredentials represents error for invalid auth credentials
	ErrInvalidCredentials = SejastipError{
		Message:    "Email atau password salah",
		ErrorCode:  401,
		HTTPStatus: http.StatusUnauthorized,
	}

	// ErrUnauthorized represents
	ErrUnauthorized = SejastipError{
		Message:    "Unauthorized",
		ErrorCode:  401,
		HTTPStatus: http.StatusUnauthorized,
	}

	// ErrForbidden represents
	ErrForbidden = SejastipError{
		Message:    "Forbidden to access the resource",
		ErrorCode:  403,
		HTTPStatus: http.StatusForbidden,
	}

	// ErrEditProductForbidden represents
	ErrEditProductForbidden = SejastipError{
		Message:    "Kamu tidak bisa mengubah atau menghapus produk yang bukan milik kamu",
		ErrorCode:  403,
		HTTPStatus: http.StatusForbidden,
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

func ValidationError(err error) SejastipError {
	return SejastipError{
		Message:    err.Error(),
		ErrorCode:  400,
		HTTPStatus: http.StatusBadRequest,
	}
}
