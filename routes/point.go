package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zihaolam/roadwarrior-backend/handlers"
)

func RegisterPointRoutes(app *fiber.App) *fiber.Router {
	pointAPIRoute := app.Group("/point")
	pointHandler := handlers.NewPointHandler()

	pointAPIRoute.Get("/:pointId", pointHandler.GetOneHandler)
	pointAPIRoute.Get("/", pointHandler.GetAllHandler)
	pointAPIRoute.Post("/", pointHandler.CreateHandler)
	pointAPIRoute.Put("/:pointId", pointHandler.UpdateHandler)
	pointAPIRoute.Delete("/:pointId", pointHandler.DeleteHandler)

	return &pointAPIRoute
}
