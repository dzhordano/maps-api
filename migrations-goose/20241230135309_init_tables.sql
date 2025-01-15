-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- Тип транспорта на маршруте
CREATE TYPE ROUTE_VEHICLE_TYPE_ENUM AS ENUM (
  'bus', -- Автобусы
  'minibus', -- Маршрутки
  'trolleybus', -- Троллейбусы
  'train' -- Электрички/Поезда
);

-- Тип маршрута
CREATE TYPE ROUTE_TYPE_ENUM AS ENUM (
  'city', -- Городской транспорт
  'intercity' -- Междугородний транспорт
);

-- Остановки
CREATE TABLE IF NOT EXISTS waypoints (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  latitude NUMERIC NOT NULL,
  longitude NUMERIC NOT NULL,
  geom geometry(Point, 4326),
  UNIQUE(latitude, longitude)
);

CREATE INDEX idx_waypoints_geom ON waypoints USING GIST(geom);

CREATE OR REPLACE FUNCTION St_DistanceSphere(lon NUMERIC, lat NUMERIC)
RETURNS TABLE(id UUID, name VARCHAR(255), latitude NUMERIC, longitude NUMERIC) AS $$
BEGIN
  RETURN QUERY
  SELECT id, name, latitude, longitude
  FROM waypoints  
  ORDER BY ST_DistanceSphere(
      geom,
      ST_SetSRID(ST_MakePoint(lon, lat), 4326)
  )
  LIMIT 1;
END;
$$ LANGUAGE plpgsql;

-- Таблица маршрутов
CREATE TABLE IF NOT EXISTS routes (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  route_kind INTEGER NOT NULL, -- Тип направления маршрута: 1 или 2
  length INTEGER NOT NULL, -- Длина маршрута (количество остановок)
  price INTEGER NOT NULL, -- Цена проезда на маршруте
  vehicle_type ROUTE_VEHICLE_TYPE_ENUM,
  route_type ROUTE_TYPE_ENUM,
  UNIQUE(name, route_kind)
);

CREATE INDEX idx_routes_name ON routes(name);

-- Связующая таблица для точек и маршрутов
CREATE TABLE IF NOT EXISTS waypoint_routes (
  waypoint_id UUID REFERENCES waypoints(id) ON DELETE CASCADE,
  route_id UUID REFERENCES routes(id) ON DELETE CASCADE,
  route_name VARCHAR(255) NOT NULL,
  route_kind INTEGER NOT NULL,
  route_number INTEGER NOT NULL,
  PRIMARY KEY(waypoint_id, route_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS waypoint_routes;
DROP TABLE IF EXISTS routes;
DROP TABLE IF EXISTS waypoints;
DROP TYPE IF EXISTS ROUTE_VEHICLE_TYPE_ENUM;
DROP TYPE IF EXISTS ROUTE_TYPE_ENUM;
DROP INDEX IF EXISTS idx_routes_name;
DROP INDEX IF EXISTS idx_waypoints_geom;
-- +goose StatementEnd
