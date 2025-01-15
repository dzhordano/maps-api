package route

import (
	"net/http"

	"github.com/dzhordano/maps-api/internal/delivery/http/v1/controller"
	"github.com/dzhordano/maps-api/internal/delivery/http/v1/route/middleware"
	"github.com/dzhordano/maps-api/pkg/logger"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupV1(log logger.Logger, wc *controller.WaypointsController, rc *controller.RoutesController, r *chi.Mux) {
	v1 := chi.NewRouter()

	// Инициализация пути получения файла с документацией
	r.Get("/docs/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/openapi.yaml")
	})

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:9000/docs/swagger.yaml")))

	r.Mount("/api/v1", v1)

	v1.Group(func(r chi.Router) {
		r.Use(middleware.CORS())

		NewWaypointsRouter(log, wc, r)
		NewRoutesRouter(log, rc, r)
	})
}
