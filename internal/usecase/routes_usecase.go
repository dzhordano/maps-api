package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/dzhordano/maps-api/internal/domain"
	"github.com/dzhordano/maps-api/pkg/logger"
	"github.com/google/uuid"
)

type routesUsecase struct {
	repo domain.RoutesRepository
	log  logger.Logger
}

func NewRoutesUsecase(repo domain.RoutesRepository, log logger.Logger) domain.RoutesUsecase {
	return &routesUsecase{
		repo: repo,
		log:  log,
	}
}

func (r *routesUsecase) List(ctx context.Context, limit, offset uint64) ([]domain.Route, error) {
	routes, err := r.repo.List(ctx, limit, offset)
	if err != nil {

		r.log.Error("list routes", "error:", err)

		return nil, domain.ErrInternalServerError
	}

	return routes, nil
}

func (r *routesUsecase) GetById(ctx context.Context, id uuid.UUID) (domain.Route, []domain.Waypoint, error) {
	route, err := r.repo.GetById(ctx, id)
	if err != nil {

		r.log.Error("get route by id", "error:", err)

		if errors.Is(err, domain.ErrNotFound) {
			return domain.Route{}, nil, fmt.Errorf("%w: route not found", domain.ErrNotFound)
		}

		return domain.Route{}, nil, domain.ErrInternalServerError
	}

	waypointIds, err := r.repo.RouteWaypoints(ctx, id, route.RouteKind)
	if err != nil {

		r.log.Error("get route by id", "error:", err)

		return domain.Route{}, nil, domain.ErrInternalServerError
	}

	return route, waypointIds, nil
}

func (r *routesUsecase) Create(ctx context.Context, route domain.Route, waypointIds []uuid.UUID) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	route.ID = id

	if err = r.repo.Create(ctx, route, waypointIds); err != nil {

		r.log.Error("create route", "error:", err)

		if errors.Is(err, domain.ErrConflict) {
			return fmt.Errorf("%w: route already exists", domain.ErrConflict)
		}

		if errors.Is(err, domain.ErrBadRequest) {
			return fmt.Errorf("%w: invalid waypoints request", domain.ErrBadRequest)
		}

		return domain.ErrInternalServerError
	}

	return nil
}

func (r *routesUsecase) Update(ctx context.Context, route domain.Route) error {
	panic("implement me")
}

func (r *routesUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.repo.Delete(ctx, id); err != nil {

		r.log.Error("delete route", "error:", err)

		if errors.Is(err, domain.ErrNotFound) {
			return fmt.Errorf("%w: route not found", domain.ErrNotFound)
		}

		return domain.ErrInternalServerError
	}

	return nil
}
