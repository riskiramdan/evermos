package product

import (
	"context"
	"errors"
	"time"

	"github.com/riskiramdan/evermos/internal/appcontext"
	"github.com/riskiramdan/evermos/internal/data"
	"github.com/riskiramdan/evermos/internal/types"
)

// Errors
var (
	// ErrWrongPassword      = errors.New("wrong password")
	// ErrWrongEmail         = errors.New("wrong email")
	ErrProductAlreadyExists = errors.New("Product Already Exists")
	ErrProductUnavailable   = errors.New("Product Unavailable")

// ErrNotFound           = errors.New("not found")
// ErrNoInput            = errors.New("no input")
// ErrLimitInput         = errors.New("Name should be more than 5 char")
// ErrNameAlreadyExist   = errors.New(("Name Already Exits"))
)

// Product product
type Product struct {
	ID        int        `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Qty       int        `json:"qty" db:"qty"`
	Price     int        `json:"price" db:"price"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
}

// OrderHistory OrderHistory
type OrderHistory struct {
	ID        int        `json:"id" db:"id"`
	UserID    int        `json:"userId" db:"user_id"`
	ProductID int        `json:"productId" db:"product_id"`
	Qty       int        `json:"qty" db:"qty"`
	Price     int        `json:"price" db:"price"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
}

//FindAllProductsParams params for find all
type FindAllProductsParams struct {
	Page      int    `json:"page"`
	Search    string `json:"search"`
	Limit     int    `json:"limit"`
	ProductID int    `json:"productId"`
	Name      string `json:"name"`
}

//FindAllOrderHistorysParams params for find all
type FindAllOrderHistorysParams struct {
	Page   int    `json:"page"`
	Search string `json:"search"`
	Limit  int    `json:"limit"`
	ID     int    `json:"id"`
}

// TransactionProductParams represent the http request data for create product
type TransactionProductParams struct {
	Name  string `json:"name"`
	Qty   int    `json:"qty"`
	Price int    `json:"price"`
}

// TransactionOrderHistorytParams represent the http request data for create order history
type TransactionOrderHistorytParams struct {
	ProductID int `json:"productId"`
	Qty       int `json:"qty"`
	Price     int `json:"price"`
}

// Storage represents the product storage interface
type Storage interface {
	FindAll(ctx context.Context, params *FindAllProductsParams) ([]*Product, *types.Error)
	FindByID(ctx context.Context, productID int) (*Product, *types.Error)
	Insert(ctx context.Context, product *Product) (*Product, *types.Error)
	Update(ctx context.Context, product *Product) (*Product, *types.Error)
	Delete(ctx context.Context, productID int) *types.Error
}

// StorageOrderHistory represents the order History storage interface
type StorageOrderHistory interface {
	FindAllOrderHistory(ctx context.Context, params *FindAllOrderHistorysParams) ([]*OrderHistory, *types.Error)
	FindOrderHistoryByID(ctx context.Context, orderHistoryID int) (*OrderHistory, *types.Error)
	InsertOrderHistory(ctx context.Context, OrderHistory *OrderHistory) (*OrderHistory, *types.Error)
	UpdateOrderHistory(ctx context.Context, OrderHistory *OrderHistory) (*OrderHistory, *types.Error)
	DeleteOrderHistory(ctx context.Context, orderHistoryID int) *types.Error
}

// ServiceInterface represents the product service interface
type ServiceInterface interface {
	ListProducts(ctx context.Context, params *FindAllProductsParams) ([]*Product, int, *types.Error)
	GetProduct(ctx context.Context, productID int) (*Product, *types.Error)
	CreateProduct(ctx context.Context, params *TransactionProductParams) (*Product, *types.Error)
	UpdateProduct(ctx context.Context, productID int, params *TransactionProductParams) (*Product, *types.Error)
	DeleteProduct(ctx context.Context, productID int) *types.Error
	CreateOrder(ctx context.Context, params *TransactionOrderHistorytParams) (*OrderHistory, *types.Error)
}

// Service is the domain logic implementation of product Service interface
type Service struct {
	productStorage Storage
	orderStorage   StorageOrderHistory
}

// ListProducts is listing products
func (s *Service) ListProducts(ctx context.Context, params *FindAllProductsParams) ([]*Product, int, *types.Error) {
	products, err := s.productStorage.FindAll(ctx, params)
	if err != nil {
		err.Path = ".ProductService->ListProducts()" + err.Path
		return nil, 0, err
	}
	params.Page = 0
	params.Limit = 0
	allProducts, err := s.productStorage.FindAll(ctx, params)
	if err != nil {
		err.Path = ".ProductService->ListProducts()" + err.Path
		return nil, 0, err
	}

	return products, len(allProducts), nil
}

// GetProduct is get product
func (s *Service) GetProduct(ctx context.Context, productID int) (*Product, *types.Error) {
	product, err := s.productStorage.FindByID(ctx, productID)
	if err != nil {
		err.Path = ".ProductService->GetProduct()" + err.Path
		return nil, err
	}

	return product, nil
}

// CreateProduct create product
func (s *Service) CreateProduct(ctx context.Context, params *TransactionProductParams) (*Product, *types.Error) {
	products, _, errType := s.ListProducts(ctx, &FindAllProductsParams{
		Name: params.Name,
	})
	if errType != nil {
		errType.Path = ".ProductService->CreateProduct()" + errType.Path
		return nil, errType
	}
	if len(products) > 0 {
		return nil, &types.Error{
			Path:    ".ProductService->CreateProduct()",
			Message: ErrProductAlreadyExists.Error(),
			Error:   ErrProductAlreadyExists,
			Type:    "validation-error",
		}
	}

	now := time.Now()

	product := &Product{
		Name:      params.Name,
		Qty:       params.Qty,
		Price:     params.Price,
		CreatedAt: now,
		UpdatedAt: &now,
	}

	product, errType = s.productStorage.Insert(ctx, product)
	if errType != nil {
		errType.Path = ".ProductService->CreateProduct()" + errType.Path
		return nil, errType
	}

	return product, nil
}

// UpdateProduct update a product
func (s *Service) UpdateProduct(ctx context.Context, productID int, params *TransactionProductParams) (*Product, *types.Error) {
	product, err := s.GetProduct(ctx, productID)
	if err != nil {
		err.Path = ".ProductService->UpdateProduct()" + err.Path
		return nil, err
	}

	products, _, err := s.ListProducts(ctx, &FindAllProductsParams{
		Name: params.Name,
	})
	if err != nil {
		err.Path = ".ProductService->UpdateProduct()" + err.Path
		return nil, err
	}
	if params.Name != "" {
		if len(products) > 0 {
			return nil, &types.Error{
				Path:    ".ProductService->CreateProduct()",
				Message: data.ErrAlreadyExist.Error(),
				Error:   data.ErrAlreadyExist,
				Type:    "validation-error",
			}
		}
	}
	if params.Name != "" {
		product.Name = params.Name
	}

	product.Qty = params.Qty
	product.Price = params.Price

	product, err = s.productStorage.Update(ctx, product)
	if err != nil {
		err.Path = ".ProductService->UpdateProduct()" + err.Path
		return nil, err
	}

	return product, nil
}

// DeleteProduct delete a product
func (s *Service) DeleteProduct(ctx context.Context, productID int) *types.Error {
	err := s.productStorage.Delete(ctx, productID)
	if err != nil {
		err.Path = ".ProductService->DeleteProduct()" + err.Path
		return err
	}

	return nil
}

// CreateOrder for create order
func (s *Service) CreateOrder(ctx context.Context, params *TransactionOrderHistorytParams) (*OrderHistory, *types.Error) {
	product, errType := s.GetProduct(ctx, params.ProductID)
	if errType != nil {
		errType.Path = ".ProductService->CreateOrder()" + errType.Path
		return nil, errType
	}
	if product.Qty < params.Qty {
		return nil, &types.Error{
			Path:    ".ProductService->CreateOrder()",
			Message: ErrProductUnavailable.Error(),
			Error:   ErrProductUnavailable,
			Type:    "validation-error",
		}
	}

	product, errType = s.UpdateProduct(ctx, params.ProductID, &TransactionProductParams{
		Qty: product.Qty - params.Qty,
	})
	if errType != nil {
		errType.Path = ".ProductService->CreateOrder()" + errType.Path
		return nil, errType
	}

	/*

		Payment Logic Here

	*/

	now := time.Now()

	orderHistory, errType := s.orderStorage.InsertOrderHistory(ctx, &OrderHistory{
		ProductID: params.ProductID,
		UserID:    appcontext.UserID(ctx),
		Qty:       params.Qty,
		Price:     params.Price,
		CreatedAt: now,
		UpdatedAt: &now,
	})

	if errType != nil {
		errType.Path = ".ProductService->CreateProduct()" + errType.Path
		return nil, errType
	}

	return orderHistory, nil
}

// NewService creates a new product AppService
func NewService(
	productStorage Storage,
	orderStorage StorageOrderHistory,
) *Service {
	return &Service{
		productStorage: productStorage,
		orderStorage:   orderStorage,
	}
}
