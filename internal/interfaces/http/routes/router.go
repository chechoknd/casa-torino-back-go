package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/casatorino/backend/internal/interfaces/http/handlers"
	appmiddleware "github.com/casatorino/backend/internal/interfaces/http/middleware"
)

type Dependencies struct {
	Auth        *handlers.AuthHandler
	Customers   *handlers.CustomerHandler
	Products    *handlers.ProductHandler
	Ingredients *handlers.IngredientHandler
	Recipes     *handlers.RecipeHandler
	Orders      *handlers.OrderHandler
	Payments    *handlers.PaymentHandler

	CORSAllowedOrigins []string
	TokenVerifier      appmiddleware.TokenVerifier
}

func NewRouter(deps Dependencies) http.Handler {
	router := chi.NewRouter()
	router.Use(appmiddleware.RequestID)
	router.Use(appmiddleware.CORS(deps.CORSAllowedOrigins))
	router.Use(appmiddleware.ContentType)
	router.Use(appmiddleware.Logger)
	router.Use(appmiddleware.Recoverer)

	if deps.Auth != nil {
		router.Route("/auth", func(r chi.Router) {
			r.Post("/register", deps.Auth.Register)
			r.Post("/login", deps.Auth.Login)
		})
	}

	protected := chi.NewRouter()
	if deps.TokenVerifier != nil {
		protected.Use(appmiddleware.JWTAuth(deps.TokenVerifier))
	}

	protected.Route("/customers", func(r chi.Router) {
		r.Post("/", deps.Customers.Create)
		r.Get("/", deps.Customers.List)
		r.Get("/{id}", deps.Customers.Get)
		r.Put("/{id}", deps.Customers.Update)
		r.Delete("/{id}", deps.Customers.Delete)
	})

	protected.Route("/products", func(r chi.Router) {
		r.Post("/", deps.Products.Create)
		r.Get("/", deps.Products.List)
		r.Get("/{id}", deps.Products.Get)
		r.Put("/{id}", deps.Products.Update)
		r.Delete("/{id}", deps.Products.Delete)
	})

	protected.Route("/ingredients", func(r chi.Router) {
		r.Post("/", deps.Ingredients.Create)
		r.Get("/", deps.Ingredients.List)
		r.Get("/{id}", deps.Ingredients.Get)
		r.Put("/{id}", deps.Ingredients.Update)
		r.Delete("/{id}", deps.Ingredients.Delete)
	})

	protected.Route("/recipes", func(r chi.Router) {
		r.Post("/", deps.Recipes.Create)
		r.Get("/", deps.Recipes.List)
		r.Post("/{id}/items", deps.Recipes.AddItem)
		r.Get("/product/{product_id}", deps.Recipes.GetByProduct)
		r.Get("/{id}/cost", deps.Recipes.GetCost)
	})

	protected.Route("/orders", func(r chi.Router) {
		r.Post("/", deps.Orders.Create)
		r.Get("/", deps.Orders.List)
		r.Get("/{id}", deps.Orders.Get)
		r.Post("/{id}/items", deps.Orders.AddItem)
		r.Patch("/{id}/status", deps.Orders.UpdateStatus)
		r.Get("/{id}/payments", deps.Payments.GetByOrder)
	})

	protected.Route("/payments", func(r chi.Router) {
		r.Get("/", deps.Payments.List)
		r.Post("/", deps.Payments.Create)
		r.Patch("/{id}/status", deps.Payments.UpdateStatus)
	})

	router.Mount("/", protected)

	return router
}
