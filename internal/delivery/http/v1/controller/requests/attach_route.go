package requests

import (
	"errors"

	"github.com/google/uuid"
)

type AttachRouteRequest struct {
	RouteID     string `json:"route_id"`
	RouteKind   int    `json:"route_kind"`
	RouteNumber int    `json:"route_number"`
}

func (r *AttachRouteRequest) Validate() error {
	if _, err := uuid.Parse(r.RouteID); err != nil {
		return errors.New("invalid route id")
	}

	if r.RouteNumber < 0 {
		return errors.New("invalid route number")
	}

	return nil
}
