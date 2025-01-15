package requests

import (
	"errors"
	"strings"

	"github.com/dzhordano/maps-api/internal/domain"
)

type UpdateRouteRequest struct {
	Name        *string `json:"name"`
	Price       *int    `json:"price"`
	VehicleType *string `json:"vehicle_type"`
	RouteType   *string `json:"route_type"`
}

func (r UpdateRouteRequest) Validate() error {
	if r.Name != nil {
		if strings.TrimSpace(*r.Name) == "" || len(*r.Name) > 256 {
			return errors.New("invalid name")
		}
	}

	if r.Price != nil {
		if *r.Price < 0 {
			return errors.New("invalid price")
		}
	}

	if r.VehicleType != nil {
		if strings.TrimSpace(*r.VehicleType) == "" || len(*r.VehicleType) > 256 || !domain.ValidVehicleType(*r.VehicleType) {
			return errors.New("invalid vehicle type")
		}
	}

	if r.RouteType != nil {
		if strings.TrimSpace(*r.RouteType) == "" || len(*r.RouteType) > 256 || !domain.ValidRouteType(*r.RouteType) {
			return errors.New("invalid route type")
		}
	}

	return nil

}
