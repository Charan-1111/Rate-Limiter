package server

import (
	"fmt"
	"goapp/store"
	"goapp/utils"

	"github.com/redis/go-redis/v9"
)

type Application struct {
	config *utils.Config
	rdb    *redis.Client
}

func NewApplication(filePath string) (*Application, error) {
	// Load the configuration file
	config := &utils.Config{}
	if err := config.LoadConfig(filePath); err != nil {
		return nil, err
	}

	// Initialize Redis
	rdb := store.InitRedis(&config.Redis)

	return &Application{
		config: config,
		rdb:    rdb,
	}, nil
}

func (app *Application) StartServer() error {
	// Start fiber server
	app.StartFiberServer()

	return nil
}

func (app *Application) StartFiberServer()  {
	appServer := app.SetupRoutes()

	if err := appServer.Listen(app.config.Ports.FiberServer); err != nil {
		fmt.Println("Error starting fiber server:", err)
	}
}