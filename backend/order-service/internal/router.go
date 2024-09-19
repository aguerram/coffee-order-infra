package internal

import (
	"github.com/aguerram/coffee-order-app/order-service/internal/handler"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	orderHandler := handler.NewOrderHandler(db)

	orderGroup := app.Group("/orders")
	orderGroup.Post("/", orderHandler.CreateOrder)
}
