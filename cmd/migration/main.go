package main

import (
	"github.com/dzhordano/maps-api/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.MustNew()

	db, err := goose.OpenDBWithDriver("postgres", cfg.PG.DSN)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	if err := goose.Up(db, "migrations-goose"); err != nil {
		panic(err)
	}
}
