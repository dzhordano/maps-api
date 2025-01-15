run:
	docker-compose up --remove-orphans

include .env
run.db:
	docker run --name=maps-db -p ${POSTGRES_PORT}:5432 -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgis/postgis

run.go: run.db run.migration
	go build -o .bin/app ./cmd/app/main.go
	./.bin/app

run.migration:
	sleep 5
	go run cmd/migration/main.go

stop.db:
	docker rm -f maps-db

# enter database container
exec.db:
	docker exec -it maps-db psql -U postgres
