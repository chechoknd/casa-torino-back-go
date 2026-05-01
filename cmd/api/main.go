package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	customeruc "github.com/casatorino/backend/internal/application/usecases/customer"
	ingredientuc "github.com/casatorino/backend/internal/application/usecases/ingredient"
	orderuc "github.com/casatorino/backend/internal/application/usecases/order"
	paymentuc "github.com/casatorino/backend/internal/application/usecases/payment"
	productuc "github.com/casatorino/backend/internal/application/usecases/product"
	recipeuc "github.com/casatorino/backend/internal/application/usecases/recipe"
	"github.com/casatorino/backend/internal/infrastructure/config"
	"github.com/casatorino/backend/internal/infrastructure/database/postgres"
	"github.com/casatorino/backend/internal/interfaces/http/handlers"
	"github.com/casatorino/backend/internal/interfaces/http/routes"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conn, err := postgres.NewConnection(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect to postgres: %v", err)
	}
	defer conn.Close()

	customerRepo := postgres.NewCustomerRepository(conn)
	productRepo := postgres.NewProductRepository(conn)
	ingredientRepo := postgres.NewIngredientRepository(conn)
	recipeRepo := postgres.NewRecipeRepository(conn)
	orderRepo := postgres.NewOrderRepository(conn)
	paymentRepo := postgres.NewPaymentRepository(conn)

	customerUseCase := customeruc.NewUseCase(customerRepo)
	productUseCase := productuc.NewUseCase(productRepo)
	ingredientUseCase := ingredientuc.NewUseCase(ingredientRepo)
	recipeUseCase := recipeuc.NewUseCase(recipeRepo, productRepo, ingredientRepo)
	orderUseCase := orderuc.NewUseCase(orderRepo, customerRepo, productRepo)
	paymentUseCase := paymentuc.NewUseCase(paymentRepo, orderRepo)

	router := routes.NewRouter(routes.Dependencies{
		Health:      handlers.NewHealthHandler(),
		Customers:   handlers.NewCustomerHandler(customerUseCase),
		Products:    handlers.NewProductHandler(productUseCase),
		Ingredients: handlers.NewIngredientHandler(ingredientUseCase),
		Recipes:     handlers.NewRecipeHandler(recipeUseCase),
		Orders:      handlers.NewOrderHandler(orderUseCase),
		Payments:    handlers.NewPaymentHandler(paymentUseCase),
	})

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("http server listening on port %s in %s mode", cfg.Port, cfg.Env)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("http server shutdown error: %v", err)
		os.Exit(1)
	}
}
