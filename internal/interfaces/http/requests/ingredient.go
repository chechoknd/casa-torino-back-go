package requests

type CreateIngredientRequest struct {
	Name         string `json:"name"`
	Unit         string `json:"unit"`
	AverageCost  string `json:"average_cost"`
	Stock        string `json:"stock"`
	MinimumStock string `json:"minimum_stock"`
}

type UpdateIngredientRequest struct {
	Name         string `json:"name"`
	Unit         string `json:"unit"`
	AverageCost  string `json:"average_cost"`
	Stock        string `json:"stock"`
	MinimumStock string `json:"minimum_stock"`
}
