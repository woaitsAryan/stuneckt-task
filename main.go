package main

import (
	"github.com/woaitsAryan/stuneckt-task/config"
	"github.com/woaitsAryan/stuneckt-task/helpers"
	"github.com/woaitsAryan/stuneckt-task/initializers"
	"github.com/woaitsAryan/stuneckt-task/routers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.AddLogger()
	initializers.ConnectToCache()
	initializers.AutoMigrate()
}

func main() {
	defer initializers.LoggerCleanUp()

	app := fiber.New(fiber.Config{
		ErrorHandler: helpers.ErrorHandler,
	})

	app.Use(helmet.New())
	app.Use(config.CORS())

	app.Use(logger.New())

	app.Static("/", "./public")

	routers.Config(app)

	app.Listen(":" + initializers.CONFIG.PORT)
}
