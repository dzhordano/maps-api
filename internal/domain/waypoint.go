package domain

import (
	"context"

	"github.com/google/uuid"
)

type Waypoint struct {
	ID        uuid.UUID
	Name      string
	Latitude  float64
	Longitude float64
}

type CommonRoutes struct {
	From   Waypoint
	To     Waypoint
	Routes []Route
}

type WaypointsRepository interface {
	List(ctx context.Context, limit, offset uint64) ([]Waypoint, error)
	GetById(ctx context.Context, id uuid.UUID) (Waypoint, error)
	GetOfNearest(ctx context.Context, amount int, latitude, longitude float64) ([]Waypoint, error)
	Create(ctx context.Context, waypoint Waypoint) error
	Update(ctx context.Context, waypoint Waypoint) error
	Delete(ctx context.Context, id uuid.UUID) error

	// TODO Вынести в отдельный интерфейс
	ListRoutes(ctx context.Context, wID uuid.UUID) ([]WaypointRoute, error)
	WaypointRoutes(ctx context.Context, wID uuid.UUID) ([]WaypointRoute, error)
}

type WaypointsUsecase interface {
	List(ctx context.Context, limit, offset uint64) ([]Waypoint, error)
	GetById(ctx context.Context, id uuid.UUID) (Waypoint, error)
	GetOfNearest(ctx context.Context, amount int, latitude, longitude float64) ([]Waypoint, error)
	Create(ctx context.Context, waypoint Waypoint) error
	Update(ctx context.Context, waypoint Waypoint) error
	Delete(ctx context.Context, id uuid.UUID) error

	CollectRoutes(ctx context.Context, amount int, lat1, long1, lat2, long2 float64) ([]CommonRoutes, error)

	// TODO Вынести в отдельный интерфейс
	ListRoutes(ctx context.Context, wID uuid.UUID) ([]WaypointRoute, error)
	CommonRoutes(ctx context.Context, w1, w2 uuid.UUID) ([]Route, error)
}
