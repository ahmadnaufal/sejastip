package api

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// ResponseBody is our default structure for API responses
type ResponseBody struct {
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Meta    interface{} `json:"meta"`
}

// MetaInfo defines additional data to be embedded in response
type MetaInfo struct {
	Status int `json:"status"`
}

// MetaPagination is an extended version of MetaInfo with pagination info
type MetaPagination struct {
	Status int `json:"status"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

// NewMetaPagination to create pagination meta
func NewMetaPagination(status, limit, offset, total int) MetaPagination {
	return MetaPagination{
		Status: status,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}
}

type ErrorBody struct {
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// OK is a wrapper to return 200 OK responses
func OK(w http.ResponseWriter, data interface{}, msg string) {
	response := ResponseBody{
		Data:    data,
		Message: msg,
		Meta:    MetaInfo{http.StatusOK},
	}
	write(w, response, http.StatusOK)
}

// OKWithMeta is a wrapper to return 200 OK responses with customized metadata
func OKWithMeta(w http.ResponseWriter, data interface{}, msg string, meta interface{}) {
	response := ResponseBody{
		Data:    data,
		Message: msg,
		Meta:    meta,
	}
	write(w, response, http.StatusOK)
}

// Created is a wrapper to return 201 Created responses
func Created(w http.ResponseWriter, data interface{}, msg string) {
	response := ResponseBody{
		Data:    data,
		Message: msg,
		Meta:    MetaInfo{http.StatusCreated},
	}
	write(w, response, http.StatusCreated)
}

func Error(w http.ResponseWriter, err error) {
	var errBody ErrorBody
	status := http.StatusInternalServerError

	switch origin := errors.Cause(err).(type) {
	case SejastipError:
		errBody = ErrorBody{
			Message: origin.Message,
			Code:    origin.ErrorCode,
		}
		status = origin.HTTPStatus
	default:
		errBody = ErrorBody{
			Message: "Internal Server Error",
			Code:    999,
		}
	}

	response := ResponseBody{
		Error: &errBody,
		Meta:  MetaInfo{status},
	}
	write(w, response, status)
}

func write(w http.ResponseWriter, result interface{}, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(result)
}
