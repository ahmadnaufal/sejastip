package api

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// ResponseBody is our default structure for API responses
type ResponseBody struct {
	Data  interface{} `json:"data,omitempty"`
	Error *ErrorBody  `json:"error,omitempty"`
	Meta  interface{} `json:"meta"`
}

type ErrorBody struct {
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// OK is a wrapper to return 200 OK responses
func OK(w http.ResponseWriter, data interface{}, meta interface{}) {
	response := ResponseBody{
		Data: data,
		Meta: meta,
	}
	write(w, response, http.StatusOK)
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
	}
	write(w, response, status)
}

func write(w http.ResponseWriter, result interface{}, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(result)
}
