package api

import (
	"net/http"
	"net/url"
	"strconv"
)

type PaginationMeta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

// NewPaginationMeta to create pagination meta
func NewPaginationMeta(limit, offset, total int) PaginationMeta {
	return PaginationMeta{
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}
}

// QueryHelper represent helper to get query url data
type QueryHelper struct {
	r  *http.Request
	uv url.Values
}

// GetString to get query url value with string data type, return empty string if query url not found
func (q *QueryHelper) GetString(p string, vars ...string) string {
	defValue := ""
	if len(vars) > 0 {
		defValue = vars[0]
	}

	sv := q.uv.Get(p)
	if sv != "" {
		return sv
	}
	return defValue
}

// GetInt to get query url value with integer data type, return 0 if query url not found
func (q *QueryHelper) GetInt(p string, vars ...int) int {
	defValue := 0
	if len(vars) > 0 {
		defValue = vars[0]
	}
	sv := q.uv.Get(p)
	if sv != "" {
		if v, err := strconv.Atoi(sv); err == nil {
			return v
		}
	}
	return defValue
}

// GetFloat to get query url value with float data type, return 0.0 if query url not found
func (q *QueryHelper) GetFloat(p string, vars ...float64) float64 {
	defValue := 0.0
	if len(vars) > 0 {
		defValue = vars[0]
	}
	sv := q.uv.Get(p)
	if sv != "" {
		if v, err := strconv.ParseFloat(sv, 64); err == nil {
			return v
		}
	}
	return defValue
}

// GetBool to get query url value with boolean data type, return false if query url not found
func (q *QueryHelper) GetBool(p string, vars ...bool) bool {
	defValue := false
	if len(vars) > 0 {
		defValue = vars[0]
	}
	sv := q.uv.Get(p)
	if sv != "" {
		v, err := strconv.ParseBool(sv)
		if err != nil {
			return defValue
		}
		return v
	}
	return defValue
}

// NewQueryHelper is a function to create query helper struct
func NewQueryHelper(r *http.Request) *QueryHelper {
	return &QueryHelper{r, r.URL.Query()}
}
