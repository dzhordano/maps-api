package requests

import (
	"errors"
	"fmt"
)

type CreateWaypointRequest struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

func (r CreateWaypointRequest) Validate() error {
	fmt.Println("waypoint", r)

	if r.Name == "" || len(r.Name) > 256 {
		return errors.New("invalid name")
	}

	if r.Latitude < -90 || r.Latitude > 90 {
		return errors.New("invalid lat")
	}

	if r.Longitude < -180 || r.Longitude > 180 {
		return errors.New("invalid lon")
	}

	return nil
}
