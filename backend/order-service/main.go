package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/aguerram/coffee-order-app/order-service/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto Migrate the Coffee model
	err = internal.AutoMigrate(db)
	if err != nil {
		panic("failed to perform auto migration")
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
		if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
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
