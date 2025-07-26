package main

import (

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/pranavpatil6/go_mart/database"
	"github.com/pranavpatil6/go_mart/routes"
)
func main() {
	godotenv.Load()

	database.ConnectDb()

	app := fiber.New()

	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Go-Mart!")
	})

	routes.SetupRoutes(app)

	app.Listen(":3000")
}