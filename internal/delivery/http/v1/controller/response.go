package controller

import (
	"encoding/json"
	"net/http"

	"github.com/dzhordano/maps-api/internal/domain"
	"github.com/pkg/errors"
)

type respErr struct {
	Error string `json:"error"`
}

func httpResponse(w http.ResponseWriter, code int, err string) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(respErr{Error: err})
}

func DomainErrorToHTTP(err error) int {
	switch errors.Unwrap(err) {
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	case domain.ErrBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
