package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/zihaolam/roadwarrior-backend/handlers"
	"github.com/zihaolam/roadwarrior-backend/repos"
	"github.com/zihaolam/roadwarrior-backend/routes"

	"github.com/gofiber/fiber/v2"
)

func migrate() {
	repos.PointRepo.Migrate()
}

func main() {
	migrate()
	// Make the handler available for Remote Procedure Call by AWS Lambda
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Pong")
	})
	app.Static("/", "./web-build")
	app.Post("/migrate", func(c *fiber.Ctx) error {
		points, err := repos.PointRepo.GetAll()
		if err != nil {
			return handlers.ErrorResponse(c, fiber.StatusInternalServerError, err)
		}

		err = repos.PointRepo.Replicate(points)
		if err != nil {
			return handlers.ErrorResponse(c, fiber.StatusInternalServerError, err)
		}
		return handlers.OKResponse(c, points)
	})

	routes.RegisterDocsRoute(app)
	routes.RegisterPointRoutes(app)

	if err := app.Listen(":80"); err != nil {
		fmt.Println(err.Error())
	}
}
