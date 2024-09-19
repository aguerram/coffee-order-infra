package internal

import (
	"github.com/aguerram/coffee-order-app/order-service/internal/model"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.Order{})
}
