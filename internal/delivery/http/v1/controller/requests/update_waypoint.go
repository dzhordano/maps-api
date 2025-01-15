package requests

import "errors"

type UpdateWaypointRequest struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

func (r UpdateWaypointRequest) Validate() error {
	if r.Name == "" || len(r.Name) > 100 {
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
