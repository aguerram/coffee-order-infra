package handler

import (
	"log"
	"math"
	"math/rand/v2"
	"strconv"

	"github.com/aguerram/coffee-order-app/inventory-service/internal/model"
	"github.com/go-faker/faker/v4"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CoffeeHandler struct {
	db *gorm.DB
}

func NewCoffeeHandler(db *gorm.DB) *CoffeeHandler {
	return &CoffeeHandler{db: db}
}

func (h CoffeeHandler) GetCoffeesHandler(c *fiber.Ctx) error {
	log.Printf("Attempt to get coffees")
	coffees := []model.Coffee{}
	tx := h.db.Find(&coffees)
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "error",
			"message": "Failed to get coffees",
			"data":    []string{},
		})
	}
	log.Printf("Coffees fetched successfully: %v", len(coffees))

	coffeesDTO := make([]model.CoffeeDTO, len(coffees))
	for i, coffee := range coffees {
		coffeesDTO[i] = coffee.ToDTO()
	}
	return c.JSON(&fiber.Map{
		"status":  "success",
		"message": "Coffees fetched successfully",
		"data":    coffeesDTO,
	})
}
func (h CoffeeHandler) PostRandomCoffeeHandler(c *fiber.Ctx) error {
	price := 0.5 + rand.Float64()*(10-0.5)
	coffee := model.Coffee{
		Title: faker.Name(),
		Price: math.Round(price*100) / 100, // Round to 2 decimal places
	}

	log.Printf("Attempt to create coffee: %v", coffee)

	tx := h.db.Create(&coffee)
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "error",
			"message": "Failed to create coffee",
			"data":    []string{},
		})
	}

	return c.JSON(&fiber.Map{
		"status":  "success",
		"message": "Coffee created successfully",
		"data":    []string{},
	})
}

func (h CoffeeHandler) GetCoffeeByIdHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	// Validate that the id parameter is numeric
	if _, err := strconv.Atoi(id); err != nil {
		log.Printf("Invalid coffee ID: %v", id)
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "error",
			"message": "Invalid coffee ID. Must be a number.",
			"data":    []string{},
		})
	}
	coffee := model.Coffee{}
	tx := h.db.First(&coffee, id)
	if tx.Error != nil {
		log.Printf("Failed to get coffee: %v", tx.Error)
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"status":  "error",
			"message": "Coffee not found",
			"data":    []string{},
		})
	}
	log.Printf("Coffee fetched successfully: %v", coffee)
	coffeeDTO := coffee.ToDTO()
	return c.JSON(&fiber.Map{
		"status":  "success",
		"message": "Coffee fetched successfully",
		"data":    coffeeDTO,
	})
}
