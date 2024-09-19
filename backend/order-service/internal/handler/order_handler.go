package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aguerram/coffee-order-app/order-service/internal/model"
	"github.com/aguerram/coffee-order-app/order-service/internal/request"
	"github.com/aguerram/coffee-order-app/order-service/internal/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type OrderHandler struct {
	db *gorm.DB
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{db: db}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var request request.CreateOrderRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorResponse := make(map[string]string)
		for _, e := range validationErrors {
			errorResponse[e.Field()] = e.Tag()
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Validation failed",
			"errors":  errorResponse,
		})
	}

	inventoryServiceUrl := os.Getenv("INVENTORY_SERVICE_URL")

	inventoryResponse, err := http.Get(fmt.Sprintf("%s/coffees/%d", inventoryServiceUrl, request.CoffeeId))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer inventoryResponse.Body.Close()

	if inventoryResponse.StatusCode == http.StatusNotFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("Coffee with ID %d does not exist", request.CoffeeId),
		})
	}

	inventoryData, err := io.ReadAll(inventoryResponse.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	inventory := &response.InventoryApiResponse{}
	err = json.Unmarshal(inventoryData, inventory)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	log.Printf("Inventory: %v", inventory.Data)
	order := model.Order{
		CoffeeId:   request.CoffeeId,
		Quantity:   request.Quantity,
		TotalPrice: float64(request.Quantity) * inventory.Data.Price,
		Status:     "pending",
	}

	tx := h.db.Create(&order)
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": tx.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(order.ToDTO())
}
