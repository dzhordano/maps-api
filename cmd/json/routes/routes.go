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
	rs := jsonbus.GetRoutes("json/myroutes.json")

	fmt.Println("routes", rs)

	routes := make([]reqmodels.CreateRouteRequest, len(rs))

	for i, r := range rs {
		routes[i] = reqmodels.CreateRouteRequest{
			Name:        r.Name,
			RouteKind:   r.RouteKind,
			Length:      r.Length,
			Price:       r.Price,
			VehicleType: r.VehicleType,
			RouteType:   r.RouteType,
			Waypoints:   r.Waypoints,
		}
	}

	c := &http.Client{}

	for _, r := range routes {
		j, err := json.Marshal(r)
		if err != nil {
			panic(err)
		}

		req, err := http.NewRequest("POST", "http://localhost:9000/api/v1/routes", bytes.NewBuffer(j))
		if err != nil {
			panic(err)
		}

		r, err := c.Do(req)
		if err != nil {
			panic(err)
		}

		if r.StatusCode != 201 {
			fmt.Println("error:", r.StatusCode)
			fmt.Println("wp:", r)
		}
	}
}
