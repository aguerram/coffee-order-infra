package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/aguerram/coffee-order-app/order-service/internal"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Env struct {
	DSN  string `validate:"required"`
	PORT string `validate:"required"`
}

func main() {
	godotenv.Load()

	validate := validator.New()
	env := Env{
		DSN:  os.Getenv("DSN"),
		PORT: os.Getenv("PORT"),
	}

	err := validate.Struct(env)
	if err != nil {
		log.Panicf("failed to validate environment variables: %v", err)
	}

	db, err := gorm.Open(postgres.Open(env.DSN), &gorm.Config{})
	if err != nil {
		log.Panicf("failed to connect database: %v", err)
	}

	// Auto Migrate the Coffee model
	err = internal.AutoMigrate(db)
	if err != nil {
		log.Panicf("failed to perform auto migration: %v", err)
	}
	fmt.Println("Database connection established and migrations completed")

	app := fiber.New()

	internal.RegisterFiberCorsMiddleware(app)

	internal.RegisterRoutes(app, db)

	// Create a channel to listen for interrupt signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Start the server in a goroutine
	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", env.PORT)); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	// Block until we receive an interrupt signal
	<-c

	// Attempt to gracefully shutdown the server
	fmt.Println("Gracefully shutting down...")
	if err := app.Shutdown(); err != nil {
		fmt.Printf("Error during shutdown: %v\n", err)
	}
}
