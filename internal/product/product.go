package product

import (
	"context"
	"errors"
	"time"

	"github.com/riskiramdan/evermos/internal/data"
	"github.com/riskiramdan/evermos/internal/types"
)

// Errors
var (
	// ErrWrongPassword      = errors.New("wrong password")
	// ErrWrongEmail         = errors.New("wrong email")
	ErrProductAlreadyExists = errors.New("Product Already Exists")

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

//FindAllProductsParams params for find all
type FindAllProductsParams struct {
	Page      int    `json:"page"`
	Search    string `json:"search"`
	Limit     int    `json:"limit"`
	ProductID int    `json:"productId"`
	Name      string `json:"name"`
}

// TransactionProductParams represent the http request data for create product
// swagger:model
type TransactionProductParams struct {
	Name  string `json:"name"`
	Qty   int    `json:"qty"`
	Price int    `json:"price"`
}

// Storage represents the product storage interface
type Storage interface {
	FindAll(ctx context.Context, params *FindAllProductsParams) ([]*Product, *types.Error)
	FindByID(ctx context.Context, productID int) (*Product, *types.Error)
	Insert(ctx context.Context, product *Product) (*Product, *types.Error)
	Update(ctx context.Context, product *Product) (*Product, *types.Error)
	Delete(ctx context.Context, productID int) *types.Error
}

// ServiceInterface represents the product service interface
type ServiceInterface interface {
	ListProducts(ctx context.Context, params *FindAllProductsParams) ([]*Product, int, *types.Error)
	GetProduct(ctx context.Context, productID int) (*Product, *types.Error)
	CreateProduct(ctx context.Context, params *TransactionProductParams) (*Product, *types.Error)
	UpdateProduct(ctx context.Context, productID int, params *TransactionProductParams) (*Product, *types.Error)
	DeleteProduct(ctx context.Context, productID int) *types.Error
}

// Service is the domain logic implementation of product Service interface
type Service struct {
	productStorage Storage
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

	if params.Qty != 0 {
		product.Qty = params.Qty
	}

	if params.Price != 0 {
		product.Price = params.Price
	}

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

// NewService creates a new product AppService
func NewService(
	productStorage Storage,
) *Service {
	return &Service{
		productStorage: productStorage,
	}
}
