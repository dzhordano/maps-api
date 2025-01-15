package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/dzhordano/maps-api/internal/domain"
	"github.com/dzhordano/maps-api/pkg/logger"
	"github.com/google/uuid"
)

type waypointsUsecase struct {
	wRepo domain.WaypointsRepository
	rRepo domain.RoutesRepository

	log logger.Logger
}

func NewWaypointsUsecase(wRepo domain.WaypointsRepository, rRepo domain.RoutesRepository, log logger.Logger) domain.WaypointsUsecase {
	return &waypointsUsecase{
		wRepo: wRepo,
		rRepo: rRepo,
		log:   log,
	}
}

func (w *waypointsUsecase) Create(ctx context.Context, waypoint domain.Waypoint) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	waypoint.ID = id

	if err := w.wRepo.Create(ctx, waypoint); err != nil {

		w.log.Error("create waypoint", "error:", err)

		if errors.Is(err, domain.ErrConflict) {
			return fmt.Errorf("%w: waypoint already exists", domain.ErrConflict)
		}

		return domain.ErrInternalServerError
	}

	return nil
}

func (w *waypointsUsecase) List(ctx context.Context, limit, offset uint64) ([]domain.Waypoint, error) {
	return w.wRepo.List(ctx, limit, offset)
}

func (w *waypointsUsecase) GetById(ctx context.Context, id uuid.UUID) (domain.Waypoint, error) {
	wp, err := w.wRepo.GetById(ctx, id)
	if err != nil {

		w.log.Error("get waypoint by id", "error:", err)

		if errors.Is(err, domain.ErrNotFound) {
			return domain.Waypoint{}, fmt.Errorf("%w: waypoint not found", domain.ErrNotFound)
		}

		return domain.Waypoint{}, domain.ErrInternalServerError
	}

	return wp, nil
}

func (w *waypointsUsecase) GetOfNearest(ctx context.Context, amount int, latitude, longitude float64) ([]domain.Waypoint, error) {
	wp, err := w.wRepo.GetOfNearest(ctx, amount, latitude, longitude)
	if err != nil {

		w.log.Error("get of nearest", "error:", err)

		return nil, domain.ErrInternalServerError
	}

	return wp, nil
}

func (w *waypointsUsecase) Update(ctx context.Context, waypoint domain.Waypoint) error {
	if err := w.wRepo.Update(ctx, waypoint); err != nil {

		w.log.Error("update waypoint", "error:", err)

		if errors.Is(err, domain.ErrNotFound) {
			return fmt.Errorf("%w: waypoint not found", domain.ErrNotFound)
		}

		return domain.ErrInternalServerError
	}

	return nil
}

func (w *waypointsUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := w.wRepo.Delete(ctx, id); err != nil {

		w.log.Error("delete waypoint", "error:", err)

		if errors.Is(err, domain.ErrNotFound) {
			return fmt.Errorf("%w: waypoint not found", domain.ErrNotFound)
		}

		return domain.ErrInternalServerError
	}

	return nil
}

func (w *waypointsUsecase) ListRoutes(ctx context.Context, wID uuid.UUID) ([]domain.WaypointRoute, error) {
	routes, err := w.wRepo.ListRoutes(ctx, wID)
	if err != nil {

		w.log.Error("list routes", "error:", err)

		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("%w: waypoint not found", domain.ErrNotFound)
		}

		return nil, domain.ErrInternalServerError
	}

	return routes, nil
}

func (w *waypointsUsecase) CommonRoutes(ctx context.Context, w1, w2 uuid.UUID) ([]domain.Route, error) {
	wr1, err := w.wRepo.WaypointRoutes(ctx, w1)
	if err != nil {

		w.log.Error("common routes", "error:", err)

		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("%w: waypoint 1 not found", domain.ErrNotFound)
		}

		return nil, domain.ErrInternalServerError
	}

	wr2, err := w.wRepo.WaypointRoutes(ctx, w2)
	if err != nil {

		w.log.Error("common routes", "error:", err)

		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("%w: waypoint 2 not found", domain.ErrNotFound)
		}

		return nil, domain.ErrInternalServerError
	}

	w.log.Debug("common routes", "w1 routes:", wr1, "w2 routes:", wr2)

	var commonRoutes []uuid.UUID
	for _, r := range wr1 {
		for _, r2 := range wr2 {
			if r.RouteID == r2.RouteID && r.RouteKind == r2.RouteKind && r.RouteNumber < r2.RouteNumber {
				commonRoutes = append(commonRoutes, r.RouteID)
			}
		}
	}

	routes, err := w.rRepo.GetByIds(ctx, commonRoutes...)
	if err != nil {

		w.log.Error("common routes", "error:", err)

		return nil, domain.ErrInternalServerError
	}

	return routes, nil
}

func (w *waypointsUsecase) CollectRoutes(ctx context.Context, waypointsAmount int, lat1, lon1, lat2, lon2 float64) ([]domain.CommonRoutes, error) {
	ws1, err := w.wRepo.GetOfNearest(ctx, waypointsAmount, lat1, lon1)
	if err != nil {

		w.log.Error("collect routes", "error:", err)

		return nil, domain.ErrInternalServerError
	}

	ws2, err := w.wRepo.GetOfNearest(ctx, waypointsAmount, lat2, lon2)
	if err != nil {

		w.log.Error("collect routes", "error:", err)

		return nil, domain.ErrInternalServerError
	}

	// var commonRoutes []domain.CommonRoutes

	for _, ws1 := range ws1 {
		for _, ws2 := range ws2 {

			routes, err := w.CommonRoutes(ctx, ws1.ID, ws2.ID)
			if err != nil {

				w.log.Error("collect routes", "error:", err)

				return nil, domain.ErrInternalServerError
			}

			if len(routes) > 0 {

				// commonRoutes = append(commonRoutes, domain.CommonRoutes{
				// 	From:   ws1,
				// 	To:     ws2,
				// 	Routes: routes,
				// })

				return append(
					[]domain.CommonRoutes{}, domain.CommonRoutes{
						From:   ws1,
						To:     ws2,
						Routes: routes,
					}), nil
			}
		}
	}

	return nil, domain.ErrNotFound
}
