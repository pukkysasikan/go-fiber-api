package routes

import (
	"fiber-mongo-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func MenuRoute(app *fiber.App) {
	app.Post("/menu", controllers.CreateMenu)
	app.Get("/menu/:menuId", controllers.GetAMenu)
	app.Put("/menu/:menuId", controllers.EditAMenu)
	app.Delete("/menu/:menuId", controllers.DeleteAMenu)
	app.Get("/menus", controllers.GetAllMenus)
}
