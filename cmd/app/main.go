package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dzhordano/maps-api/internal/config"
	"github.com/dzhordano/maps-api/internal/delivery/http/v1/controller"
	"github.com/dzhordano/maps-api/internal/delivery/http/v1/route"
	"github.com/dzhordano/maps-api/internal/httpserver"
	"github.com/dzhordano/maps-api/internal/repository"
	"github.com/dzhordano/maps-api/internal/usecase"
	"github.com/dzhordano/maps-api/pkg/databases/pg"
	"github.com/dzhordano/maps-api/pkg/logger"
	"github.com/go-chi/chi/v5"
)

// TODO

// (?) - not sure what/how.
// [?] - certainly, but skill issue for now.

// Add metrics.
// Add rate limiter.
// Add tests.

// [?] Do some caching.

func main() {
	cfg := config.MustNew()

	pool, err := pg.NewClient(context.Background(), cfg.PG.DSN)
	if err != nil {
		panic(err)
	}
	defer pg.Close(pool)

	log := logger.MustNewSlogLogger(os.Stdout, cfg.LogLevel)

	r := chi.NewRouter()

	wRepo := repository.NewWaypointRepo(pool)
	rRepo := repository.NewRoutesRepo(pool)

	rUsecase := usecase.NewRoutesUsecase(rRepo, log)
	wUsecase := usecase.NewWaypointsUsecase(wRepo, rRepo, log)

	rController := controller.NewRouteController(log, rUsecase)
	wController := controller.NewWaypointsController(log, wUsecase)

	route.SetupV1(log, wController, rController, r)

	srv := httpserver.New(net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port), r)

	log.Info("server started", "port", cfg.HTTP.Port)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

		<-quit

		log.Info("server is shutting down")

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error("server shutdown failed", "error", err)
		}
	}()

	if err := srv.Run(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
