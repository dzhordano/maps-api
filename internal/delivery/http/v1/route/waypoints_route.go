package route

import (
	"github.com/dzhordano/maps-api/internal/delivery/http/v1/controller"
	"github.com/dzhordano/maps-api/pkg/logger"
	"github.com/go-chi/chi/v5"
)

func NewWaypointsRouter(log logger.Logger, wc *controller.WaypointsController, r chi.Router) {

	// Вернуть маршруты между двумя точками (в amount). Где от каждой точки до каждой другой возвращаются общие маршруты.
	r.Get("/waypoints/route", wc.CollectRoutes)

	r.Get("/waypoints", wc.List)               // Получение всех существующих точек.
	r.Get("/waypoints/{id}", wc.Get)           // Получение точки по id c подробной информацией об остановке *пока никакой такой информации нету*.
	r.Get("/waypoints/nearest", wc.GetNearest) // Получение ближайших точек от параметров (количество, широта, долгота).

	r.Post("/waypoints", wc.Create)        // Создание новой точки (остановки).
	r.Put("/waypoints/{id}", wc.Update)    // Обновление точки.
	r.Delete("/waypoints/{id}", wc.Delete) // Удаление точки.

	r.Get("/waypoints/{id}/routes", wc.ListRoutes)                    // Получение маршрутов, проходящих через остановку.
	r.Get("/waypoints/{id}/routes/{waypoint_id}", wc.GetCommonRoutes) // Получение общих маршрутов между двумя остановками.
}
