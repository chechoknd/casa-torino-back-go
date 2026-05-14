package requests

type CreateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ProductType string `json:"product_type"`
	BasePrice   string `json:"base_price"`
	CostPrice   string `json:"cost_price"`
	ImageURL    string `json:"image_url"`
	IsPublic    bool   `json:"is_public"`
}

type UpdateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ProductType string `json:"product_type"`
	BasePrice   string `json:"base_price"`
	CostPrice   string `json:"cost_price"`
	ImageURL    string `json:"image_url"`
	IsPublic    bool   `json:"is_public"`
}
