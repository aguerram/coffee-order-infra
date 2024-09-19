package internal

import (
	"github.com/aguerram/coffee-order-app/inventory-service/internal/model"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.Coffee{})
}
