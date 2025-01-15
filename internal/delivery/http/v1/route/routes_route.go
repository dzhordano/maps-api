package route

import (
	"github.com/dzhordano/maps-api/internal/delivery/http/v1/controller"
	"github.com/dzhordano/maps-api/pkg/logger"
	"github.com/go-chi/chi/v5"
)

func NewRoutesRouter(log logger.Logger, rc *controller.RoutesController, r chi.Router) {
	r.Get("/routes/{id}", rc.GetRouteById) // Получение маршрута по id. Также возврат всех остановок на маршруте.
	r.Get("/routes", rc.ListRoutes)        // Получение всех существующих маршрутов.

	r.Post("/routes", rc.CreateRoute)        // Создание маршрута.
	r.Delete("/routes/{id}", rc.DeleteRoute) // Удаление маршрута.
}
