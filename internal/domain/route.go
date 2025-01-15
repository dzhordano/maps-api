package domain

import (
	"context"

	"github.com/google/uuid"
)

type Route struct {
	ID          uuid.UUID
	Name        string // Название маршрута
	RouteKind   int    // Вид маршрута (1 или 2). Используется для определения направления маршрута
	Length      int    // Длина маршрута (количество остановок)
	Price       int    // Цена проезда на маршруте
	VehicleType string // Тип транспорта
	RouteType   string // Тип маршрута (внутригородской, межгородской)
}

type WaypointRoute struct {
	RouteID     uuid.UUID
	WaypointID  uuid.UUID
	RouteName   string
	RouteKind   int
	RouteNumber int
}

type RoutesRepository interface {
	List(ctx context.Context, limit, offset uint64) ([]Route, error)
	GetById(ctx context.Context, id uuid.UUID) (Route, error)

	Create(ctx context.Context, route Route, waypointIds []uuid.UUID) error
	Update(ctx context.Context, route Route) error
	Delete(ctx context.Context, id uuid.UUID) error

	GetByIds(ctx context.Context, id ...uuid.UUID) ([]Route, error)
	RouteWaypoints(ctx context.Context, rID uuid.UUID, rKind int) ([]Waypoint, error)
}

type RoutesUsecase interface {
	List(ctx context.Context, limit, offset uint64) ([]Route, error)
	GetById(ctx context.Context, id uuid.UUID) (Route, []Waypoint, error)

	Create(ctx context.Context, route Route, waypointIds []uuid.UUID) error
	Update(ctx context.Context, route Route) error
	Delete(ctx context.Context, id uuid.UUID) error
}

func ValidVehicleType(vt string) bool {
	switch vt {
	case "bus":
		return true
	case "trolleybus":
		return true
	case "train":
		return true
	default:
		return false
	}
}

func ValidRouteType(rt string) bool {
	switch rt {
	case "city":
		return true
	case "intercity":
		return true
	default:
		return false
	}
}
