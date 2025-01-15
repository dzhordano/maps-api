package responses

import (
	"github.com/dzhordano/maps-api/internal/domain"
)

type GetRouteByIdResponse struct {
	Route     domain.Route      `json:"route"`
	Waypoints []domain.Waypoint `json:"waypoints"`
}
