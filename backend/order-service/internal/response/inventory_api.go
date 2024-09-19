package response

type InventoryApiResponse struct {
	Data struct {
		ID    int     `json:"id"`
		Title string  `json:"title"`
		Price float64 `json:"price"`
	} `json:"data"`
}
