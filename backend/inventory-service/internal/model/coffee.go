package model

import "gorm.io/gorm"

type Coffee struct {
	gorm.Model
	Title string  `gorm:"not null"`
	Price float64 `gorm:"not null;type:decimal(10,2)"`
}

type CoffeeDTO struct {
	ID    uint    `json:"id"`
	Title string  `json:"title"`
	Price float64 `json:"price"`
}

func (c *Coffee) ToDTO() CoffeeDTO {
	return CoffeeDTO{
		ID:    c.ID,
		Title: c.Title,
		Price: c.Price,
	}
}
