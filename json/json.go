package json

import (
	"encoding/json"
	"os"
)

type JSONWaypoint struct {
	Name        string     `json:"name"`
	Coordinates [2]float64 `json:"coordinates"`
}

func GetWaypoints(path string) []JSONWaypoint {
	var waypoints []JSONWaypoint

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	json.NewDecoder(f).Decode(&waypoints)

	if len(waypoints) == 0 {
		panic("no waypoints")
	}

	return waypoints
}

type JSONRoute struct {
	Name        string   `json:"name"`
	RouteKind   int      `json:"route_kind"`
	Length      int      `json:"length"`
	Price       int      `json:"price"`
	VehicleType string   `json:"vehicle_type"`
	RouteType   string   `json:"route_type"`
	Waypoints   []string `json:"waypoints"`
}

func GetRoutes(path string) []JSONRoute {
	var routes []JSONRoute

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	json.NewDecoder(f).Decode(&routes)

	if len(routes) == 0 {
		panic("no routes")
	}

	return routes
}
