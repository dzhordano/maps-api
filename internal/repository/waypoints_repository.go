package repository

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/dzhordano/maps-api/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	waypointTable       = "waypoints"
	waypointRoutesTable = "waypoint_routes"
)

type waypointRepo struct {
	db *pgxpool.Pool
}

func NewWaypointRepo(db *pgxpool.Pool) domain.WaypointsRepository {
	return &waypointRepo{db: db}
}

func (r *waypointRepo) Create(ctx context.Context, waypoint domain.Waypoint) error {
	insertBuilder := sq.Insert(waypointTable).
		Columns("id", "name", "latitude", "longitude", "geom").
		Values(waypoint.ID, waypoint.Name, waypoint.Latitude, waypoint.Longitude,
			fmt.Sprintf("SRID=4326;POINT(%f %f)", waypoint.Longitude, waypoint.Latitude)).
		PlaceholderFormat(sq.Dollar)

	query, args, err := insertBuilder.ToSql()
	if err != nil {

		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	if err != nil {
		var pqErr *pgconn.PgError

		if errors.As(err, &pqErr) {
			if pqErr.Code == pgerrcode.UniqueViolation {
				return fmt.Errorf("%w, %s", domain.ErrConflict, err)
			}
		}

		return err
	}

	return nil
}

func (r *waypointRepo) GetById(ctx context.Context, id uuid.UUID) (domain.Waypoint, error) {
	selectBuilder := sq.Select("id", "name", "latitude", "longitude").
		From(waypointTable).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := selectBuilder.ToSql()
	if err != nil {

		return domain.Waypoint{}, err
	}

	var waypoint domain.Waypoint
	if err := r.db.QueryRow(ctx, query, args...).Scan(&waypoint.ID, &waypoint.Name, &waypoint.Latitude, &waypoint.Longitude); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Waypoint{}, fmt.Errorf("%w, %s", domain.ErrNotFound, err)
		}

		return domain.Waypoint{}, err
	}

	return waypoint, nil
}

func (r *waypointRepo) List(ctx context.Context, limit, offset uint64) ([]domain.Waypoint, error) {
	selectBuilder := sq.Select("id", "name", "latitude", "longitude").
		From(waypointTable).
		Limit(limit).
		Offset(offset).
		PlaceholderFormat(sq.Dollar)

	query, args, err := selectBuilder.ToSql()
	if err != nil {

		return nil, err
	}

	var waypoints []domain.Waypoint
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var waypoint domain.Waypoint
		if err := rows.Scan(&waypoint.ID, &waypoint.Name, &waypoint.Latitude, &waypoint.Longitude); err != nil {

			return nil, err
		}
		waypoints = append(waypoints, waypoint)
	}

	if err := rows.Err(); err != nil {

		return nil, err
	}

	return waypoints, nil
}

func (r *waypointRepo) GetOfNearest(ctx context.Context, amount int, latitude, longitude float64) ([]domain.Waypoint, error) {
	query := `
    SELECT id, name, latitude, longitude
    FROM waypoints
    ORDER BY ST_DistanceSphere(
        geom,
        ST_SetSRID(ST_MakePoint($1, $2), 4326)
    )
    LIMIT $3;
	`

	var waypoints []domain.Waypoint

	rows, err := r.db.Query(ctx, query, longitude, latitude, amount)
	if err != nil {

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var waypoint domain.Waypoint
		if err := rows.Scan(&waypoint.ID, &waypoint.Name, &waypoint.Latitude, &waypoint.Longitude); err != nil {

			return nil, err
		}
		waypoints = append(waypoints, waypoint)
	}

	return waypoints, nil
}

func (r *waypointRepo) Update(ctx context.Context, waypoint domain.Waypoint) error {
	updateBuilder := sq.Update(waypointTable).
		Set("name", waypoint.Name).
		Set("latitude", waypoint.Latitude).
		Set("longitude", waypoint.Longitude).
		Where(sq.Eq{"id": waypoint.ID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := updateBuilder.ToSql()
	if err != nil {

		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%w, %s", domain.ErrNotFound, err)
		}

		var pqErr *pgconn.PgError
		if errors.As(err, &pqErr) {
			if pqErr.Code == pgerrcode.UniqueViolation {
				return fmt.Errorf("%w, %s", domain.ErrConflict, err)
			}
		}

		return err
	}

	return nil
}

func (r *waypointRepo) Delete(ctx context.Context, id uuid.UUID) error {
	deleteBuilder := sq.Delete(waypointTable).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := deleteBuilder.ToSql()
	if err != nil {

		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%w, %s", domain.ErrNotFound, err)
		}

		return err
	}

	return nil
}

func (r *waypointRepo) ListRoutes(ctx context.Context, wID uuid.UUID) ([]domain.WaypointRoute, error) {
	selectBuilder := sq.Select("route_id", "route_name", "route_kind", "route_number").
		From(waypointRoutesTable).
		Where(sq.Eq{"waypoint_id": wID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := selectBuilder.ToSql()
	if err != nil {

		return nil, err
	}

	var routes []domain.WaypointRoute
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w, %s", domain.ErrNotFound, err)
		}

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var route domain.WaypointRoute
		route.WaypointID = wID
		if err := rows.Scan(&route.RouteID, &route.RouteName, &route.RouteKind, &route.RouteNumber); err != nil {

			return nil, err
		}
		routes = append(routes, route)
	}

	return routes, nil
}

func (r *waypointRepo) WaypointRoutes(ctx context.Context, wID uuid.UUID) ([]domain.WaypointRoute, error) {
	selectBuilder := sq.Select("route_id", "route_number").
		From(waypointRoutesTable).
		Where(sq.Eq{"waypoint_id": wID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := selectBuilder.ToSql()
	if err != nil {

		return nil, err
	}

	var routes []domain.WaypointRoute
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w, %s", domain.ErrNotFound, err)
		}

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var route domain.WaypointRoute
		if err := rows.Scan(&route.RouteID, &route.RouteNumber); err != nil {

			return nil, err
		}
		routes = append(routes, route)
	}

	return routes, nil
}
