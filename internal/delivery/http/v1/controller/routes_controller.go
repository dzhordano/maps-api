package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dzhordano/maps-api/internal/delivery/http/v1/controller/mapper"
	"github.com/dzhordano/maps-api/internal/delivery/http/v1/controller/requests"
	"github.com/dzhordano/maps-api/internal/delivery/http/v1/controller/responses"
	"github.com/dzhordano/maps-api/internal/domain"
	"github.com/dzhordano/maps-api/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const (
	defaultLimitValue  = 10
	defaultOffsetValue = 0
)

type RoutesController struct {
	Log          logger.Logger
	RouteUsecase domain.RoutesUsecase
}

func NewRouteController(log logger.Logger, routeUsecase domain.RoutesUsecase) *RoutesController {
	return &RoutesController{
		Log:          log,
		RouteUsecase: routeUsecase,
	}
}

func (rc *RoutesController) ListRoutes(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	limit := r.URL.Query().Get("limit")

	var limitInt uint64
	var err error

	if limit == "" {
		limitInt = defaultLimitValue
	} else {
		limitInt, err = parseUint64(limit)
		if err != nil {
			httpResponse(w, http.StatusBadRequest, "invalid limit param")
			return
		}
	}

	offset := r.URL.Query().Get("offset")

	var offsetInt uint64

	if offset == "" {
		offsetInt = defaultOffsetValue
	} else {
		offsetInt, err = parseUint64(offset)
		if err != nil {
			httpResponse(w, http.StatusBadRequest, "invalid offset param")
			return
		}
	}

	routes, err := rc.RouteUsecase.List(r.Context(), limitInt, offsetInt)
	if err != nil {
		httpResponse(w, DomainErrorToHTTP(err), err.Error())
		return
	}

	rc.Log.Debug("list routes", "routes:", routes)

	err = json.NewEncoder(w).Encode(routes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rc *RoutesController) GetRouteById(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	parsedId, err := uuid.Parse(id)
	if err != nil {
		httpResponse(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	rc.Log.Debug("get route by id", "parsed id:", id)

	route, routeWaypoints, err := rc.RouteUsecase.GetById(r.Context(), parsedId)
	if err != nil {
		httpResponse(w, DomainErrorToHTTP(err), err.Error())
		return
	}

	rc.Log.Debug("get route by id", "route:", route)

	err = json.NewEncoder(w).Encode(
		responses.GetRouteByIdResponse{
			Route:     route,
			Waypoints: routeWaypoints,
		},
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rc *RoutesController) CreateRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var route requests.CreateRouteRequest
	err := json.NewDecoder(r.Body).Decode(&route)
	if err != nil {
		httpResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := route.Validate(); err != nil {
		httpResponse(w, http.StatusBadRequest, err.Error()) // TODO Perhaps return all the wrong fields
		return
	}

	rc.Log.Debug("create route", "decoded route:", route)

	var waypointIds []uuid.UUID

	for _, waypoint := range route.Waypoints {
		parsedWaypointId, err := uuid.Parse(waypoint)
		if err != nil {
			httpResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid waypoint id: %s", waypoint))
			return
		}

		waypointIds = append(waypointIds, parsedWaypointId)
	}

	rc.Log.Debug("create route", "waypoint ids:", waypointIds)

	err = rc.RouteUsecase.Create(r.Context(), mapper.CreateRouteRequestToDomain(route), waypointIds)
	if err != nil {
		httpResponse(w, DomainErrorToHTTP(err), err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (rc *RoutesController) UpdateRoute(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusNotImplemented)

	// w.Header().Add("Content-Type", "application/json")

	// id := chi.URLParam(r, "id")

	// parsedId, err := uuid.Parse(id)
	// if err != nil {
	// 	rc.Log.Error("update route", "error:", err)
	// 	httpResponse(w, http.StatusBadRequest, "invalid uuid")
	// 	return
	// }

	// var route requests.UpdateRouteRequest
	// err = json.NewDecoder(r.Body).Decode(&route)
	// if err != nil {
	// 	rc.Log.Error("update route", "error:", err)
	// 	httpResponse(w, http.StatusBadRequest, "invalid request body")
	// 	return
	// }

	// if err := route.Validate(); err != nil {
	// 	rc.Log.Error("create route", "error:", err)
	// 	httpResponse(w, http.StatusBadRequest, err.Error())
	// 	return
	// }

	// rc.Log.Debug("update route", "decoded route:", route)

	// domainRoute := domain.Route{}

	// domainRoute.ID = parsedId

	// if route.Name != nil {
	// 	domainRoute.Name = *route.Name
	// }
	// if route.Price != nil {
	// 	domainRoute.Price = *route.Price
	// }
	// if route.VehicleType != nil {
	// 	domainRoute.VehicleType = *route.VehicleType
	// }
	// if route.RouteType != nil {
	// 	domainRoute.RouteType = *route.RouteType
	// }

	// err = rc.RouteUsecase.Update(r.Context(), domainRoute)
	// if err != nil {
	// 	rc.Log.Error("update route", "error:", err)
	// 	httpResponse(w, DomainErrorToHTTP(err), err.Error())
	// 	return
	// }

	// w.WriteHeader(http.StatusOK)
}

func (rc *RoutesController) DeleteRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id := chi.URLParam(r, "id")

	parsedId, err := uuid.Parse(id)
	if err != nil {
		httpResponse(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	rc.Log.Debug("delete route", "parsed id:", id)

	err = rc.RouteUsecase.Delete(r.Context(), parsedId)
	if err != nil {
		httpResponse(w, DomainErrorToHTTP(err), err.Error())
		return
	}

	rc.Log.Debug("route deleted", "route id:", id)

	w.WriteHeader(http.StatusOK)
}

func parseFloat(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}

func parseInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func parseUint64(s string) (uint64, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}
