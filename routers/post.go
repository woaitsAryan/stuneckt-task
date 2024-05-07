package routers

import (
	"github.com/woaitsAryan/stuneckt-task/controllers"
	"github.com/woaitsAryan/stuneckt-task/middlewares"
	"github.com/gofiber/fiber/v2"
)

func PostRouter(app *fiber.App) {
	postRoutes := app.Group("/posts", middlewares.Protect)
	postRoutes.Post("/", controllers.AddPost)
	postRoutes.Get("/me", controllers.GetMyPosts)
	postRoutes.Get("/me/likes", controllers.GetMyLikedPosts)
	postRoutes.Get("/:postID", controllers.GetPost)
	postRoutes.Patch("/:postID", controllers.UpdatePost)
	postRoutes.Delete("/:postID", controllers.DeletePost)

	postRoutes.Get("/like/:postID", controllers.LikePost)

}