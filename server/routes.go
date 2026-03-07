package server

import "github.com/gofiber/fiber/v2"

func (app *Application) SetupRoutes() {
	appServer := fiber.New()

	// Defining the routes
	appServer.Get("")
}
