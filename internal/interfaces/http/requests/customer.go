package requests

type CreateCustomerRequest struct {
	FullName     string `json:"full_name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	CustomerType string `json:"customer_type"`
}

type UpdateCustomerRequest struct {
	FullName     string `json:"full_name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	CustomerType string `json:"customer_type"`
}
