package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/casatorino/backend/internal/interfaces/http/handlers"
	appmiddleware "github.com/casatorino/backend/internal/interfaces/http/middleware"
)

type Dependencies struct {
	Health      *handlers.HealthHandler
	Customers   *handlers.CustomerHandler
	Products    *handlers.ProductHandler
	Ingredients *handlers.IngredientHandler
	Recipes     *handlers.RecipeHandler
	Orders      *handlers.OrderHandler
	Payments    *handlers.PaymentHandler
}

func NewRouter(deps Dependencies) http.Handler {
	router := chi.NewRouter()
	router.Use(appmiddleware.RequestID)
	router.Use(appmiddleware.ContentType)
	router.Use(appmiddleware.Logger)
	router.Use(appmiddleware.Recoverer)

	router.Get("/health", deps.Health.Check)
	router.Post("/health", deps.Health.Check)

	router.Route("/customers", func(r chi.Router) {
		r.Post("/", deps.Customers.Create)
		r.Get("/", deps.Customers.List)
		r.Get("/{id}", deps.Customers.Get)
		r.Put("/{id}", deps.Customers.Update)
		r.Delete("/{id}", deps.Customers.Delete)
	})

	router.Route("/products", func(r chi.Router) {
		r.Post("/", deps.Products.Create)
		r.Get("/", deps.Products.List)
		r.Get("/{id}", deps.Products.Get)
		r.Put("/{id}", deps.Products.Update)
		r.Delete("/{id}", deps.Products.Delete)
	})

	router.Route("/ingredients", func(r chi.Router) {
		r.Post("/", deps.Ingredients.Create)
		r.Get("/", deps.Ingredients.List)
		r.Get("/{id}", deps.Ingredients.Get)
		r.Put("/{id}", deps.Ingredients.Update)
		r.Delete("/{id}", deps.Ingredients.Delete)
	})

	router.Route("/recipes", func(r chi.Router) {
		r.Post("/", deps.Recipes.Create)
		r.Post("/{id}/items", deps.Recipes.AddItem)
		r.Get("/product/{product_id}", deps.Recipes.GetByProduct)
		r.Get("/{id}/cost", deps.Recipes.GetCost)
	})

	router.Route("/orders", func(r chi.Router) {
		r.Post("/", deps.Orders.Create)
		r.Get("/", deps.Orders.List)
		r.Get("/{id}", deps.Orders.Get)
		r.Post("/{id}/items", deps.Orders.AddItem)
		r.Patch("/{id}/status", deps.Orders.UpdateStatus)
		r.Get("/{id}/payments", deps.Payments.GetByOrder)
	})

	router.Route("/payments", func(r chi.Router) {
		r.Post("/", deps.Payments.Create)
		r.Patch("/{id}/status", deps.Payments.UpdateStatus)
	})

	return router
}
