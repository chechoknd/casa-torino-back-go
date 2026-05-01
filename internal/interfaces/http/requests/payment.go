package requests

type CreatePaymentRequest struct {
	OrderID string `json:"order_id"`
	Amount  string `json:"amount"`
	Method  string `json:"method"`
	Status  string `json:"status"`
}

type UpdatePaymentStatusRequest struct {
	Status string `json:"status"`
}
