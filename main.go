package main

import (
	"fiber-mongo-api/configs"
	"fiber-mongo-api/routes" //add this

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(
		cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type, Accept",
		},
	))

	//run database
	configs.ConnectDB()

	//routes
	routes.MenuRoute(app) //add this
	routes.CategoryRoute(app)

	app.Listen(":8080")
}
