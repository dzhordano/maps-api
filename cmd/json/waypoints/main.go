package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	reqmodels "github.com/dzhordano/maps-api/internal/delivery/http/v1/controller/requests"
	jsonbus "github.com/dzhordano/maps-api/json"
)

func main() {
	wps := jsonbus.GetWaypoints("json/44.json")

	waypoints := make([]reqmodels.CreateWaypointRequest, len(wps))

	for i, wp := range wps {
		waypoints[i] = reqmodels.CreateWaypointRequest{
			Name:      wp.Name,
			Latitude:  wp.Coordinates[1],
			Longitude: wp.Coordinates[0],
		}
	}

	c := &http.Client{}

	for _, wp := range waypoints {
		j, err := json.Marshal(wp)
		if err != nil {
			panic(err)
		}

		req, err := http.NewRequest("POST", "http://localhost:9000/api/v1/waypoints", bytes.NewBuffer(j))
		if err != nil {
			panic(err)
		}

		r, err := c.Do(req)
		if err != nil {
			panic(err)
		}

		if r.StatusCode != 201 {
			fmt.Println("error:", r.StatusCode)
			fmt.Println("wp:", wp)
		}
	}
}
