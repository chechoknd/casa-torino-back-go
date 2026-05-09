package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/domain/entities"
)

type CustomerRepository struct {
	CreateFn      func(context.Context, *entities.Customer) error
	UpdateFn      func(context.Context, *entities.Customer) error
	DeactivateFn  func(context.Context, uuid.UUID, time.Time) error
	FindByIDFn    func(context.Context, uuid.UUID) (*entities.Customer, error)
	FindByEmailFn func(context.Context, string) (*entities.Customer, error)
	ListFn        func(context.Context) ([]entities.Customer, error)
}

type UserRepository struct {
	CreateFn         func(context.Context, *entities.User) error
	FindByIDFn       func(context.Context, uuid.UUID) (*entities.User, error)
	FindByEmailFn    func(context.Context, string) (*entities.User, error)
	FindByUsernameFn func(context.Context, string) (*entities.User, error)
}

type RefreshTokenRepository struct {
	CreateFn          func(context.Context, *entities.RefreshToken) error
	FindByTokenHashFn func(context.Context, string) (*entities.RefreshToken, error)
	RevokeFn          func(context.Context, uuid.UUID, time.Time) error
}

func (m *UserRepository) Create(ctx context.Context, user *entities.User) error {
	return m.CreateFn(ctx, user)
}
func (m *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return m.FindByIDFn(ctx, id)
}
func (m *UserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	return m.FindByEmailFn(ctx, email)
}
func (m *UserRepository) FindByUsername(ctx context.Context, username string) (*entities.User, error) {
	return m.FindByUsernameFn(ctx, username)
}

func (m *RefreshTokenRepository) Create(ctx context.Context, token *entities.RefreshToken) error {
	return m.CreateFn(ctx, token)
}
func (m *RefreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*entities.RefreshToken, error) {
	return m.FindByTokenHashFn(ctx, tokenHash)
}
func (m *RefreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID, revokedAt time.Time) error {
	return m.RevokeFn(ctx, id, revokedAt)
}

func (m *CustomerRepository) Create(ctx context.Context, customer *entities.Customer) error {
	return m.CreateFn(ctx, customer)
}
func (m *CustomerRepository) Update(ctx context.Context, customer *entities.Customer) error {
	return m.UpdateFn(ctx, customer)
}
func (m *CustomerRepository) Deactivate(ctx context.Context, id uuid.UUID, updatedAt time.Time) error {
	return m.DeactivateFn(ctx, id, updatedAt)
}
func (m *CustomerRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Customer, error) {
	return m.FindByIDFn(ctx, id)
}
func (m *CustomerRepository) FindByEmail(ctx context.Context, email string) (*entities.Customer, error) {
	return m.FindByEmailFn(ctx, email)
}
func (m *CustomerRepository) List(ctx context.Context) ([]entities.Customer, error) {
	return m.ListFn(ctx)
}

type ProductRepository struct {
	CreateFn     func(context.Context, *entities.Product) error
	UpdateFn     func(context.Context, *entities.Product) error
	DeactivateFn func(context.Context, uuid.UUID, time.Time) error
	FindByIDFn   func(context.Context, uuid.UUID) (*entities.Product, error)
	ListActiveFn func(context.Context) ([]entities.Product, error)
}

func (m *ProductRepository) Create(ctx context.Context, product *entities.Product) error {
	return m.CreateFn(ctx, product)
}
func (m *ProductRepository) Update(ctx context.Context, product *entities.Product) error {
	return m.UpdateFn(ctx, product)
}
func (m *ProductRepository) Deactivate(ctx context.Context, id uuid.UUID, updatedAt time.Time) error {
	return m.DeactivateFn(ctx, id, updatedAt)
}
func (m *ProductRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	return m.FindByIDFn(ctx, id)
}
func (m *ProductRepository) ListActive(ctx context.Context) ([]entities.Product, error) {
	return m.ListActiveFn(ctx)
}

type IngredientRepository struct {
	CreateFn     func(context.Context, *entities.Ingredient) error
	UpdateFn     func(context.Context, *entities.Ingredient) error
	DeactivateFn func(context.Context, uuid.UUID, time.Time) error
	FindByIDFn   func(context.Context, uuid.UUID) (*entities.Ingredient, error)
	ListActiveFn func(context.Context) ([]entities.Ingredient, error)
}

func (m *IngredientRepository) Create(ctx context.Context, ingredient *entities.Ingredient) error {
	return m.CreateFn(ctx, ingredient)
}
func (m *IngredientRepository) Update(ctx context.Context, ingredient *entities.Ingredient) error {
	return m.UpdateFn(ctx, ingredient)
}
func (m *IngredientRepository) Deactivate(ctx context.Context, id uuid.UUID, updatedAt time.Time) error {
	return m.DeactivateFn(ctx, id, updatedAt)
}
func (m *IngredientRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Ingredient, error) {
	return m.FindByIDFn(ctx, id)
}
func (m *IngredientRepository) ListActive(ctx context.Context) ([]entities.Ingredient, error) {
	return m.ListActiveFn(ctx)
}

type RecipeRepository struct {
	CreateFn          func(context.Context, *entities.Recipe) error
	UpdateFn          func(context.Context, *entities.Recipe) error
	AddItemFn         func(context.Context, uuid.UUID, *entities.RecipeItem) error
	FindByIDFn        func(context.Context, uuid.UUID) (*entities.Recipe, error)
	FindByProductIDFn func(context.Context, uuid.UUID) (*entities.Recipe, error)
	ListFn            func(context.Context) ([]entities.Recipe, error)
}

func (m *RecipeRepository) Create(ctx context.Context, recipe *entities.Recipe) error {
	return m.CreateFn(ctx, recipe)
}
func (m *RecipeRepository) Update(ctx context.Context, recipe *entities.Recipe) error {
	return m.UpdateFn(ctx, recipe)
}
func (m *RecipeRepository) AddItem(ctx context.Context, recipeID uuid.UUID, item *entities.RecipeItem) error {
	return m.AddItemFn(ctx, recipeID, item)
}
func (m *RecipeRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Recipe, error) {
	return m.FindByIDFn(ctx, id)
}
func (m *RecipeRepository) FindByProductID(ctx context.Context, productID uuid.UUID) (*entities.Recipe, error) {
	return m.FindByProductIDFn(ctx, productID)
}
func (m *RecipeRepository) List(ctx context.Context) ([]entities.Recipe, error) {
	return m.ListFn(ctx)
}

type OrderRepository struct {
	CreateFn           func(context.Context, *entities.Order) error
	UpdateFn           func(context.Context, *entities.Order) error
	AddItemFn          func(context.Context, uuid.UUID, *entities.OrderItem) error
	FindByIDFn         func(context.Context, uuid.UUID) (*entities.Order, error)
	ListFn             func(context.Context) ([]entities.Order, error)
	ListByCustomerIDFn func(context.Context, uuid.UUID) ([]entities.Order, error)
}

func (m *OrderRepository) Create(ctx context.Context, order *entities.Order) error {
	return m.CreateFn(ctx, order)
}
func (m *OrderRepository) Update(ctx context.Context, order *entities.Order) error {
	return m.UpdateFn(ctx, order)
}
func (m *OrderRepository) AddItem(ctx context.Context, orderID uuid.UUID, item *entities.OrderItem) error {
	return m.AddItemFn(ctx, orderID, item)
}
func (m *OrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Order, error) {
	return m.FindByIDFn(ctx, id)
}
func (m *OrderRepository) List(ctx context.Context) ([]entities.Order, error) {
	return m.ListFn(ctx)
}
func (m *OrderRepository) ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]entities.Order, error) {
	return m.ListByCustomerIDFn(ctx, customerID)
}

type PaymentRepository struct {
	CreateFn        func(context.Context, *entities.Payment) error
	UpdateFn        func(context.Context, *entities.Payment) error
	FindByIDFn      func(context.Context, uuid.UUID) (*entities.Payment, error)
	ListFn          func(context.Context) ([]entities.Payment, error)
	ListByOrderIDFn func(context.Context, uuid.UUID) ([]entities.Payment, error)
}

func (m *PaymentRepository) Create(ctx context.Context, payment *entities.Payment) error {
	return m.CreateFn(ctx, payment)
}
func (m *PaymentRepository) Update(ctx context.Context, payment *entities.Payment) error {
	return m.UpdateFn(ctx, payment)
}
func (m *PaymentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Payment, error) {
	return m.FindByIDFn(ctx, id)
}
func (m *PaymentRepository) List(ctx context.Context) ([]entities.Payment, error) {
	return m.ListFn(ctx)
}
func (m *PaymentRepository) ListByOrderID(ctx context.Context, orderID uuid.UUID) ([]entities.Payment, error) {
	return m.ListByOrderIDFn(ctx, orderID)
}
