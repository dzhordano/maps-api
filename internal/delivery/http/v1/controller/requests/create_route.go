package requests

import (
	"errors"

	"github.com/dzhordano/maps-api/internal/domain"
)

type CreateRouteRequest struct {
	Name        string   `json:"name"`
	RouteKind   int      `json:"route_kind"`
	Length      int      `json:"length"`
	Price       int      `json:"price"`
	VehicleType string   `json:"vehicle_type"`
	RouteType   string   `json:"route_type"`
	Waypoints   []string `json:"waypoints"`
}

func (r CreateRouteRequest) Validate() error {
	if r.Name == "" || len(r.Name) > 256 {
		return errors.New("invalid name")
	}

	if r.RouteKind < 0 || r.RouteKind > 2 {
		return errors.New("invalid route kind")
	}

	if r.Length < 0 {
		return errors.New("invalid length")
	}

	if r.Price < 0 {
		return errors.New("invalid price")
	}

	if r.VehicleType == "" || len(r.VehicleType) > 256 || !domain.ValidVehicleType(r.VehicleType) {
		return errors.New("invalid vehicle type")
	}

	if r.RouteType == "" || len(r.RouteType) > 256 || !domain.ValidRouteType(r.RouteType) {
		return errors.New("invalid route type")
	}

	if len(r.Waypoints) < r.Length {
		return errors.New("invalid waypoints")
	}

	return nil
}
