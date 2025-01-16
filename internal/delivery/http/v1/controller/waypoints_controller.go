package controller

import (
	"encoding/json"
	"net/http"

	"github.com/dzhordano/maps-api/internal/delivery/http/v1/controller/mapper"
	"github.com/dzhordano/maps-api/internal/delivery/http/v1/controller/requests"
	"github.com/dzhordano/maps-api/internal/domain"
	"github.com/dzhordano/maps-api/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const defaultAmountValue = 1

type WaypointsController struct {
	Log             logger.Logger
	WaypointUsecase domain.WaypointsUsecase
}

func NewWaypointsController(log logger.Logger, wu domain.WaypointsUsecase) *WaypointsController {
	return &WaypointsController{
		Log:             log,
		WaypointUsecase: wu,
	}
}

func (wc *WaypointsController) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var waypoint requests.CreateWaypointRequest

	err := json.NewDecoder(r.Body).Decode(&waypoint)
	if err != nil {
		httpResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := waypoint.Validate(); err != nil {
		httpResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	wc.Log.Debug("create waypoint", "decoded waypoint:", waypoint)

	err = wc.WaypointUsecase.Create(r.Context(), mapper.CreateWaypointRequestToDomain(waypoint))
	if err != nil {
		httpResponse(w, DomainErrorToHTTP(err), err.Error())
		return
	}

	wc.Log.Debug("waypoint created", "waypoint:", waypoint)

	w.WriteHeader(http.StatusCreated)
}

func (wc *WaypointsController) List(w http.ResponseWriter, r *http.Request) {
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
	waypoints, err := wc.WaypointUsecase.List(r.Context(), limitInt, offsetInt)
	if err != nil {
		httpResponse(w, DomainErrorToHTTP(err), err.Error())
		return
	}

	wc.Log.Debug("list waypoints", "waypoints:", waypoints)

	err = json.NewEncoder(w).Encode(waypoints)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (wc *WaypointsController) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	parsedId, err := uuid.Parse(id)
	if err != nil {
		httpResponse(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	wc.Log.Debug("get waypoint", "parsed id:", id)

	waypoint, err := wc.WaypointUsecase.GetById(r.Context(), parsedId)
	if err != nil {
		httpResponse(w, DomainErrorToHTTP(err), err.Error())
		return
	}

	wc.Log.Debug("get waypoint", "waypoint:", waypoint)

	err = json.NewEncoder(w).Encode(waypoint)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (wc *WaypointsController) GetNearest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	amount := r.URL.Query().Get("amount")
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")

	wc.Log.Debug("getOfNearest waypoint", "amount:", amount, "lat:", lat, "lon:", lon)

	if lat == "" || lon == "" {
		httpResponse(w, http.StatusBadRequest, "no lat & lon provided")
		return
	}

	latf, err := parseFloat(lat)
	if err != nil {
		httpResponse(w, http.StatusBadRequest, "invalid lat parameter")
		return
	}

	lonf, err := parseFloat(lon)
	if err != nil {
		httpResponse(w, http.StatusBadRequest, "invalid lon parameter")
		return
	}

	var amountInt int
	if amount == "" {
		amountInt = defaultAmountValue
	} else {
		amountInt, err = parseInt(amount)
		if err != nil {
			httpResponse(w, http.StatusBadRequest, "invalid amount parameter")
			return
		}
	}

	waypoints, err := wc.WaypointUsecase.GetOfNearest(r.Context(), amountInt, latf, lonf)
	if err != nil {
		httpResponse(w, DomainErrorToHTTP(err), err.Error())
		return
	}

	wc.Log.Debug("getOfNearest waypoint", "waypoints:", waypoints)

	err = json.NewEncoder(w).Encode(waypoints)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (wc *WaypointsController) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var waypoint requests.UpdateWaypointRequest

	err := json.NewDecoder(r.Body).Decode(&waypoint)
	if err != nil {
		httpResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := waypoint.Validate(); err != nil {
		httpResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	wc.Log.Debug("update waypoint", "decoded waypoint:", waypoint)

	id := chi.URLParam(r, "id")

	parsedId, err := uuid.Parse(id)
	if err != nil {
		httpResponse(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	wc.Log.Debug("update waypoint", "parsed id:", id)

	domainWaypoint := mapper.UpdateWaypointRequestToDomain(waypoint)

	domainWaypoint.ID = parsedId

	err = wc.WaypointUsecase.Update(r.Context(), domainWaypoint)
	if err != nil {
		httpResponse(w, DomainErrorToHTTP(err), err.Error())
		return
	}

	wc.Log.Debug("waypoint updated", "waypoint id:", id)

	w.WriteHeader(http.StatusOK)
}

func (wc *WaypointsController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	parsedId, err := uuid.Parse(id)
	if err != nil {
		httpResponse(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	wc.Log.Debug("delete waypoint", "parsed id:", id)

	err = wc.WaypointUsecase.Delete(r.Context(), parsedId)
	if err != nil {
		httpResponse(w, DomainErrorToHTTP(err), err.Error())
		return
	}

	wc.Log.Debug("waypoint deleted", "waypoint id:", id)

	w.WriteHeader(http.StatusOK)
}

func (wc *WaypointsController) ListRoutes(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	parsedId, err := uuid.Parse(id)
	if err != nil {
		httpResponse(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	wc.Log.Debug("list routes", "parsed id:", id)

	routes, err := wc.WaypointUsecase.ListRoutes(r.Context(), parsedId)
	if err != nil {
		httpResponse(w, DomainErrorToHTTP(err), err.Error())
		return
	}

	wc.Log.Debug("list routes", "routes:", routes)

	err = json.NewEncoder(w).Encode(routes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (wc *WaypointsController) GetCommonRoutes(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	parsedId, err := uuid.Parse(id)
	if err != nil {
		httpResponse(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	wc.Log.Debug("get common routes", "parsed id:", parsedId)

	waypointId := chi.URLParam(r, "waypoint_id")

	parsedWaypointId, err := uuid.Parse(waypointId)
	if err != nil {
		httpResponse(w, http.StatusBadRequest, "invalid waypoint uuid")
		return
	}

	wc.Log.Debug("get common routes", "parsed waypoint id:", parsedWaypointId)

	routes, err := wc.WaypointUsecase.CommonRoutes(r.Context(), parsedId, parsedWaypointId)
	if err != nil {
		httpResponse(w, DomainErrorToHTTP(err), err.Error())
		return
	}

	wc.Log.Debug("get common routes", "waypoints:", routes)

	err = json.NewEncoder(w).Encode(routes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (wc *WaypointsController) CollectRoutes(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	amount := r.URL.Query().Get("amount")

	var waypointsAmount int
	var err error

	if amount == "" {
		waypointsAmount = defaultAmountValue
	} else {
		waypointsAmount, err = parseInt(r.URL.Query().Get("amount"))
		if err != nil {
			httpResponse(w, http.StatusBadRequest, "invalid amount parameter")
			return
		}
	}

	lat1f, err := parseFloat(r.URL.Query().Get("lat1"))
	if err != nil {
		httpResponse(w, http.StatusBadRequest, "invalid lat1 parameter")
		return
	}

	lon1f, err := parseFloat(r.URL.Query().Get("lon1"))
	if err != nil {
		httpResponse(w, http.StatusBadRequest, "invalid lon1 parameter")
		return
	}

	lat2f, err := parseFloat(r.URL.Query().Get("lat2"))
	if err != nil {
		httpResponse(w, http.StatusBadRequest, "invalid lat2 parameter")
		return
	}

	lon2f, err := parseFloat(r.URL.Query().Get("lon2"))
	if err != nil {
		httpResponse(w, http.StatusBadRequest, "invalid lon2 parameter")
		return
	}

	routes, err := wc.WaypointUsecase.CollectRoutes(r.Context(), waypointsAmount, lat1f, lon1f, lat2f, lon2f)
	if err != nil {
		httpResponse(w, DomainErrorToHTTP(err), err.Error())
		return
	}

	wc.Log.Debug("collect routes", "routes:", routes)

	err = json.NewEncoder(w).Encode(routes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
