package request

type CreateOrderRequest struct {
	CoffeeId int `json:"coffee_id" validate:"required,min=1"`
	Quantity int `json:"quantity" validate:"required,min=1"`
}
