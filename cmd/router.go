package main

import "github.com/gofiber/fiber/v2"

func (app *Application) Router() *fiber.App {
	api := fiber.New()
	v1 := api.Group("/v1/api")
	{
		v1.Get("/health", app.Health)
		v1.Post("/places", app.CreatePlace)
		v1.Get("/places/:place_id", app.GetPlace)
		v1.Get("/places", app.ListPlace)
		v1.Patch("/places/:place_id", app.UpdateOne)
		v1.Delete("/places/:place_id", app.DeletePlace)
	}
	return api
}
