package mapper

import (
	"github.com/dzhordano/maps-api/internal/delivery/http/v1/controller/requests"
	"github.com/dzhordano/maps-api/internal/domain"
)

func CreateWaypointRequestToDomain(waypoint requests.CreateWaypointRequest) domain.Waypoint {
	return domain.Waypoint{
		Name:      waypoint.Name,
		Latitude:  waypoint.Latitude,
		Longitude: waypoint.Longitude,
	}
}

func UpdateWaypointRequestToDomain(waypoint requests.UpdateWaypointRequest) domain.Waypoint {
	return domain.Waypoint{
		Name:      waypoint.Name,
		Latitude:  waypoint.Latitude,
		Longitude: waypoint.Longitude,
	}
}

func CreateRouteRequestToDomain(route requests.CreateRouteRequest) domain.Route {
	return domain.Route{
		Name:        route.Name,
		RouteKind:   route.RouteKind,
		Length:      route.Length,
		Price:       route.Price,
		VehicleType: route.VehicleType,
		RouteType:   route.RouteType,
	}
}
