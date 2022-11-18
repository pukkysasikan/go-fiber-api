package routes

import (
	"fiber-mongo-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func CategoryRoute(app *fiber.App) {
	app.Post("/category", controllers.CreateCategory)
	app.Get("/category/:categoryId", controllers.GetACategory)
	app.Put("/category/:categoryId", controllers.EditACategory)
	app.Delete("/category/:categoryId", controllers.DeleteACategory)
	app.Get("/categorys", controllers.GetAllCategorys)
}
