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
	routesTable = "routes"
)

type routesRepo struct {
	db *pgxpool.Pool
}

func NewRoutesRepo(db *pgxpool.Pool) domain.RoutesRepository {
	return &routesRepo{db: db}
}

func (r *routesRepo) List(ctx context.Context, limit, offset uint64) ([]domain.Route, error) {
	selectBuilder := sq.Select("id", "name", "route_kind", "length", "price", "vehicle_type", "route_type").
		From(routesTable).
		Limit(limit).
		Offset(offset).
		PlaceholderFormat(sq.Dollar)

	query, args, err := selectBuilder.ToSql()
	if err != nil {

		return nil, err
	}

	var routes []domain.Route
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var route domain.Route
		if err := rows.Scan(&route.ID, &route.Name, &route.RouteKind, &route.Length, &route.Price, &route.VehicleType, &route.RouteType); err != nil {

			return nil, err
		}
		routes = append(routes, route)
	}

	return routes, nil
}

func (r *routesRepo) GetById(ctx context.Context, id uuid.UUID) (domain.Route, error) {
	selectBuilder := sq.Select("id", "name", "route_kind", "length", "price", "vehicle_type", "route_type").
		From(routesTable).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := selectBuilder.ToSql()
	if err != nil {

		return domain.Route{}, err
	}

	var route domain.Route
	if err := r.db.QueryRow(ctx, query, args...).Scan(&route.ID, &route.Name, &route.RouteKind, &route.Length, &route.Price, &route.VehicleType, &route.RouteType); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Route{}, fmt.Errorf("%w, %s", domain.ErrNotFound, err)
		}

		return domain.Route{}, err
	}

	return route, nil
}

func (r *routesRepo) GetByIds(ctx context.Context, id ...uuid.UUID) ([]domain.Route, error) {
	selectBuilder := sq.Select("id", "name", "route_kind", "length", "price", "vehicle_type", "route_type").
		From(routesTable).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := selectBuilder.ToSql()
	if err != nil {

		return nil, err
	}

	var routes []domain.Route
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var route domain.Route
		if err := rows.Scan(&route.ID, &route.Name, &route.RouteKind, &route.Length, &route.Price, &route.VehicleType, &route.RouteType); err != nil {

			return nil, err
		}
		routes = append(routes, route)
	}

	return routes, nil
}

func (r *routesRepo) Create(ctx context.Context, route domain.Route, waypointIds []uuid.UUID) (err error) {
	return runWithTx(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		insertBuilder := sq.Insert(routesTable).
			Columns("id", "name", "route_kind", "length", "price", "vehicle_type", "route_type").
			Values(route.ID, route.Name, route.RouteKind, route.Length, route.Price, route.VehicleType, route.RouteType).
			PlaceholderFormat(sq.Dollar)

		query, args, err := insertBuilder.ToSql()
		if err != nil {

			return err
		}

		_, err = tx.Exec(ctx, query, args...)

		if err != nil {
			var pqErr *pgconn.PgError

			if errors.As(err, &pqErr) {
				if pqErr.Code == pgerrcode.UniqueViolation {
					return fmt.Errorf("%w, %s", domain.ErrConflict, err)
				}
			}

			return err
		}

		for i, wID := range waypointIds {
			insertBuilder := sq.Insert(waypointRoutesTable).
				Columns("route_id", "waypoint_id", "route_name", "route_number", "route_kind").
				Values(route.ID, wID, route.Name, i+1, route.RouteKind).
				PlaceholderFormat(sq.Dollar)

			query, args, err := insertBuilder.ToSql()
			if err != nil {

				return err
			}

			_, err = tx.Exec(ctx, query, args...)

			if err != nil {
				var pqErr *pgconn.PgError

				if errors.As(err, &pqErr) {
					if pqErr.Code == pgerrcode.UniqueViolation {
						return fmt.Errorf("%w, %s", domain.ErrConflict, err)
					}

					if pqErr.Code == pgerrcode.ForeignKeyViolation {
						return fmt.Errorf("%w, %s", domain.ErrBadRequest, err)
					}
				}

				return err
			}
		}

		return nil
	})
}

// TODO implement
func (r *routesRepo) Update(ctx context.Context, route domain.Route) error {
	return nil
}

func (r *routesRepo) Delete(ctx context.Context, id uuid.UUID) error {
	deleteBuilder := sq.Delete(routesTable).
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

func (r *routesRepo) RouteWaypoints(ctx context.Context, rID uuid.UUID, rKind int) ([]domain.Waypoint, error) {
	selectBuilder := sq.Select("id", "name", "latitude", "longitude").
		From(waypointTable).
		Join(waypointRoutesTable + " ON id = waypoint_id").
		Where(sq.Eq{"route_id": rID, "route_kind": rKind}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := selectBuilder.ToSql()
	if err != nil {

		return nil, err
	}

	fmt.Println("query", query, args)

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

	return waypoints, nil
}

func runWithTx(ctx context.Context, db *pgxpool.Pool, fn func(ctx context.Context, tx pgx.Tx) error) error {
	tx, err := db.Begin(ctx)
	if err != nil {

		return err
	}

	err = fn(ctx, tx)
	if nil == err {
		return tx.Commit(ctx)
	}

	rErr := tx.Rollback(ctx)
	if nil != rErr {

		return errors.Join(err, rErr)
	}

	fmt.Println("HERE 2", err)

	return err
}
