package server

import (
	"goapp/handlers"

	"github.com/gofiber/fiber/v2"
)

func (app *Application) SetupRoutes() *fiber.App {
	appServer := fiber.New()

	configHandler := handlers.NewConfigHandler(app.ctx, app.config, app.log, app.db, app.rdb, app.factory, app.cache, app.cb)

	// Defining the routes
	appServer.Get("/api/v1/limiter", configHandler.GetLimiter)

	return appServer
}
