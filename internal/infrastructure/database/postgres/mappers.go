package postgres

import (
	sqlcdb "github.com/casatorino/backend/internal/infrastructure/database/sqlc"

	"github.com/casatorino/backend/internal/domain/entities"
	"github.com/casatorino/backend/internal/domain/valueobjects"
)

func mapCustomer(row sqlcdb.Customer) (entities.Customer, error) {
	customerType, err := valueobjects.NewCustomerType(row.CustomerType)
	if err != nil {
		return entities.Customer{}, err
	}

	return entities.Customer{
		ID:           row.ID,
		FullName:     row.FullName,
		Phone:        row.Phone,
		Email:        row.Email,
		CustomerType: customerType,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
		IsActive:     row.IsActive,
	}, nil
}

func mapProduct(row sqlcdb.Product) (entities.Product, error) {
	productType, err := valueobjects.NewProductType(row.ProductType)
	if err != nil {
		return entities.Product{}, err
	}

	return entities.Product{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		ProductType: productType,
		BasePrice:   row.BasePrice,
		CostPrice:   row.CostPrice,
		IsActive:    row.IsActive,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}, nil
}

func mapIngredient(row sqlcdb.Ingredient) (entities.Ingredient, error) {
	unit, err := valueobjects.NewUnit(row.Unit)
	if err != nil {
		return entities.Ingredient{}, err
	}

	return entities.Ingredient{
		ID:           row.ID,
		Name:         row.Name,
		Unit:         unit,
		AverageCost:  row.AverageCost,
		Stock:        row.Stock,
		MinimumStock: row.MinimumStock,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
		IsActive:     row.IsActive,
	}, nil
}

func mapRecipe(row sqlcdb.Recipe, items []sqlcdb.RecipeItem) (entities.Recipe, error) {
	mappedItems := make([]entities.RecipeItem, 0, len(items))
	for _, item := range items {
		unit, err := valueobjects.NewUnit(item.Unit)
		if err != nil {
			return entities.Recipe{}, err
		}

		mappedItems = append(mappedItems, entities.RecipeItem{
			ID:           item.ID,
			RecipeID:     item.RecipeID,
			IngredientID: item.IngredientID,
			Quantity:     item.Quantity,
			Unit:         unit,
		})
	}

	return entities.Recipe{
		ID:        row.ID,
		ProductID: row.ProductID,
		Name:      row.Name,
		Portions:  int(row.Portions),
		Items:     mappedItems,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		IsActive:  row.IsActive,
	}, nil
}

func mapOrder(row sqlcdb.Order, items []sqlcdb.OrderItem) (entities.Order, error) {
	status, err := valueobjects.NewOrderStatus(row.Status)
	if err != nil {
		return entities.Order{}, err
	}

	mappedItems := make([]entities.OrderItem, 0, len(items))
	for _, item := range items {
		mappedItems = append(mappedItems, entities.OrderItem{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  int(item.Quantity),
			UnitPrice: item.UnitPrice,
			Subtotal:  item.Subtotal,
		})
	}

	return entities.Order{
		ID:         row.ID,
		CustomerID: row.CustomerID,
		Status:     status,
		Items:      mappedItems,
		Subtotal:   row.Subtotal,
		Discount:   row.Discount,
		Total:      row.Total,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}, nil
}

func mapPayment(row sqlcdb.Payment) (entities.Payment, error) {
	method, err := valueobjects.NewPaymentMethod(row.Method)
	if err != nil {
		return entities.Payment{}, err
	}

	status, err := valueobjects.NewPaymentStatus(row.Status)
	if err != nil {
		return entities.Payment{}, err
	}

	return entities.Payment{
		ID:        row.ID,
		OrderID:   row.OrderID,
		Amount:    row.Amount,
		Method:    method,
		Status:    status,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}
