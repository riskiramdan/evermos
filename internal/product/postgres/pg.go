package postgres

import (
	"context"
	"fmt"

	"github.com/riskiramdan/evermos/internal/data"
	"github.com/riskiramdan/evermos/internal/product"
	"github.com/riskiramdan/evermos/internal/types"
)

// PostgresStorage implements the product storage service interface
type PostgresStorage struct {
	Storage data.GenericStorage
}

// FindAll find all products
func (s *PostgresStorage) FindAll(ctx context.Context, params *product.FindAllProductsParams) ([]*product.Product, *types.Error) {

	products := []*product.Product{}
	where := `"deleted_at" IS NULL`

	if params.ProductID != 0 {
		where += ` AND "id" = :productId`
	}
	if params.Name != "" {
		where += ` AND "name" = :name`
	}
	if params.Search != "" {
		where += ` AND "name" ILIKE :search`
	}
	if params.Page != 0 && params.Limit != 0 {
		where = fmt.Sprintf(`%s ORDER BY "id" DESC LIMIT :limit OFFSET :offset`, where)
	} else {
		where = fmt.Sprintf(`%s ORDER BY "id" DESC`, where)
	}

	err := s.Storage.Where(ctx, &products, where, map[string]interface{}{
		"productId": params.ProductID,
		"name":      params.Name,
		"limit":     params.Limit,
		"search":    "%" + params.Search + "%",
		"offset":    ((params.Page - 1) * params.Limit),
	})
	if err != nil {
		return nil, &types.Error{
			Path:    ".ProductPostgresStorage->FindAll()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return products, nil
}

// FindByID find product by its id
func (s *PostgresStorage) FindByID(ctx context.Context, productID int) (*product.Product, *types.Error) {
	products, err := s.FindAll(ctx, &product.FindAllProductsParams{
		ProductID: productID,
	})
	if err != nil {
		err.Path = ".ProductPostgresStorage->FindByID()" + err.Path
		return nil, err
	}

	if len(products) < 1 || products[0].ID != productID {
		return nil, &types.Error{
			Path:    ".ProductPostgresStorage->FindByID()",
			Message: data.ErrNotFound.Error(),
			Error:   data.ErrNotFound,
			Type:    "pq-error",
		}
	}

	return products[0], nil
}

// Insert insert product
func (s *PostgresStorage) Insert(ctx context.Context, product *product.Product) (*product.Product, *types.Error) {
	err := s.Storage.Insert(ctx, product)
	if err != nil {
		return nil, &types.Error{
			Path:    ".ProductPostgresStorage->Insert()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return product, nil
}

// Update update product
func (s *PostgresStorage) Update(ctx context.Context, product *product.Product) (*product.Product, *types.Error) {
	err := s.Storage.Update(ctx, product)
	if err != nil {
		return nil, &types.Error{
			Path:    ".ProductPostgresStorage->Update()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return product, nil
}

// Delete delete a product
func (s *PostgresStorage) Delete(ctx context.Context, productID int) *types.Error {
	err := s.Storage.Delete(ctx, productID)
	if err != nil {
		return &types.Error{
			Path:    ".ProductPostgresStorage->Delete()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return nil
}

// NewPostgresStorage creates new product repository service
func NewPostgresStorage(
	storage data.GenericStorage,
) *PostgresStorage {
	return &PostgresStorage{
		Storage: storage,
	}
}
