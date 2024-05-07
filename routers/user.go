package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/woaitsAryan/stuneckt-task/controllers"
	"github.com/woaitsAryan/stuneckt-task/middlewares"
	"github.com/woaitsAryan/stuneckt-task/validators"
)

func UserRouter(app *fiber.App) {
	app.Post("/signup", validators.UserCreateValidator, controllers.SignUp)
	app.Post("/login", controllers.LogIn)

	userRoutes := app.Group("/users", middlewares.Protect)
	userRoutes.Get("/me", controllers.GetMe)
}
