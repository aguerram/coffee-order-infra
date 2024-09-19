package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	CoffeeId   int     `gorm:"not null"`
	Quantity   int     `gorm:"not null"`
	TotalPrice float64 `gorm:"not null"`
	Status     string  `gorm:"not null"`
}

type OrderDTO struct {
	CoffeeId   int     `json:"coffee_id"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
}

func (o *Order) ToDTO() OrderDTO {
	return OrderDTO{
		CoffeeId:   o.CoffeeId,
		Quantity:   o.Quantity,
		TotalPrice: o.TotalPrice,
		Status:     o.Status,
	}
}
