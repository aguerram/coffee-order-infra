package internal

import (
	"github.com/aguerram/coffee-order-app/inventory-service/internal/handler"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	coffeeHandler := handler.NewCoffeeHandler(db)

	coffeeGroup := app.Group("/coffees")
	coffeeGroup.Get("/", coffeeHandler.GetCoffeesHandler)
	coffeeGroup.Post("/", coffeeHandler.PostRandomCoffeeHandler)
	coffeeGroup.Get("/:id", coffeeHandler.GetCoffeeByIdHandler)
}
