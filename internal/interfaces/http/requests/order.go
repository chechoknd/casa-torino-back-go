package requests

type CreateOrderRequest struct {
	CustomerID string `json:"customer_id"`
	Discount   string `json:"discount"`
}

type AddOrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}
