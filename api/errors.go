package api

import (
	"fmt"
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

	// ErrUnauthorized represents error if a user is unauthorized to request
	ErrUnauthorized = SejastipError{
		Message:    "Unauthorized",
		ErrorCode:  401,
		HTTPStatus: http.StatusUnauthorized,
	}

	// ErrForbidden represents error when a resource can't be accessed by the
	// requesting user
	ErrForbidden = SejastipError{
		Message:    "Kamu tidak dapat mengakses resource yang kamu minta",
		ErrorCode:  403,
		HTTPStatus: http.StatusForbidden,
	}

	// ErrEditProductForbidden represents error that thrown when a user tries to
	// edit a product data that is not owned by itself
	ErrEditProductForbidden = SejastipError{
		Message:    "Kamu tidak bisa mengubah atau menghapus produk yang bukan milik kamu",
		ErrorCode:  403,
		HTTPStatus: http.StatusForbidden,
	}

	// ErrEditTransactionForbidden represents error that thrown when a user tries to
	// edit a transaction data that is not owned by itself
	ErrEditTransactionForbidden = SejastipError{
		Message:    "Kamu tidak bisa mengubah transaksi dalam status saat ini atau transaksi ini bukan milik kamu",
		ErrorCode:  403,
		HTTPStatus: http.StatusForbidden,
	}

	// ErrInvalidTransactionStateTransition represents error that thrown when
	// a transaction state transition is invalid
	ErrInvalidTransactionStateTransition = SejastipError{
		Message:    "Status transaksi tidak valid",
		ErrorCode:  422,
		HTTPStatus: http.StatusUnprocessableEntity,
	}

	// ErrEditInvoiceForbidden represents error that thrown when a user tries to
	// edit an invoice data that is not owned by itself
	ErrEditInvoiceForbidden = SejastipError{
		Message:    "Kamu tidak bisa mengubah invoice dalam status saat ini atau invoice ini bukan milik kamu",
		ErrorCode:  403,
		HTTPStatus: http.StatusForbidden,
	}

	// ErrTransactionInvoiceExists represents error that thrown when a user tries to
	// create a new invoice on a transaction that already has an invoice
	ErrTransactionInvoiceExists = SejastipError{
		Message:    "Transaksi sudah memiliki invoice",
		ErrorCode:  422,
		HTTPStatus: http.StatusUnprocessableEntity,
	}

	// ErrBuyOwnProduct represents error that happens when a user trying to buy
	// its own product
	ErrBuyOwnProduct = SejastipError{
		Message:    "Kamu tidak dapat membeli produk yang kamu list sendiri",
		ErrorCode:  422,
		HTTPStatus: http.StatusUnprocessableEntity,
	}

	// ErrTransactionAddressNotOwned represents error that happens when a user tries
	// to create transaction with an address that is not owned by itself
	ErrTransactionAddressNotOwned = SejastipError{
		Message:    "Alamat tidak sesuai dengan alamat yang sudah kamu simpan",
		ErrorCode:  422,
		HTTPStatus: http.StatusUnprocessableEntity,
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

func CustomValidationError(errmsg string, args ...interface{}) SejastipError {
	return SejastipError{
		Message:    fmt.Sprintf(errmsg, args...),
		ErrorCode:  400,
		HTTPStatus: http.StatusBadRequest,
	}
}
