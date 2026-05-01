package requests

type CreateRecipeRequest struct {
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Portions  int    `json:"portions"`
}

type AddRecipeItemRequest struct {
	IngredientID string `json:"ingredient_id"`
	Quantity     string `json:"quantity"`
	Unit         string `json:"unit"`
}
