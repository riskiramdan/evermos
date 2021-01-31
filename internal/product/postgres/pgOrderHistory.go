package postgres

import (
	"context"
	"fmt"

	"github.com/riskiramdan/evermos/internal/data"
	"github.com/riskiramdan/evermos/internal/product"
	"github.com/riskiramdan/evermos/internal/types"
)

// FindAllOrderHistory find all Order History
func (s *PostgresStorage) FindAllOrderHistory(ctx context.Context, params *product.FindAllOrderHistorysParams) ([]*product.OrderHistory, *types.Error) {

	orderHistory := []*product.OrderHistory{}
	where := `"deleted_at" IS NULL`

	if params.ID != 0 {
		where += ` AND "id" = :id`
	}
	if params.Page != 0 && params.Limit != 0 {
		where = fmt.Sprintf(`%s ORDER BY "id" DESC LIMIT :limit OFFSET :offset`, where)
	} else {
		where = fmt.Sprintf(`%s ORDER BY "id" DESC`, where)
	}

	err := s.Storage.Where(ctx, &orderHistory, where, map[string]interface{}{
		"id":     params.ID,
		"limit":  params.Limit,
		"search": "%" + params.Search + "%",
		"offset": ((params.Page - 1) * params.Limit),
	})
	if err != nil {
		return nil, &types.Error{
			Path:    ".ProductPostgresStorage->FindAll()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return orderHistory, nil
}

// FindOrderHistoryByID find order history by its id
func (s *PostgresStorage) FindOrderHistoryByID(ctx context.Context, orderHistoryID int) (*product.OrderHistory, *types.Error) {
	products, err := s.FindAllOrderHistory(ctx, &product.FindAllOrderHistorysParams{
		ID: orderHistoryID,
	})
	if err != nil {
		err.Path = ".ProductPostgresStorage->FindByID()" + err.Path
		return nil, err
	}

	if len(products) < 1 || products[0].ID != orderHistoryID {
		return nil, &types.Error{
			Path:    ".ProductPostgresStorage->FindByID()",
			Message: data.ErrNotFound.Error(),
			Error:   data.ErrNotFound,
			Type:    "pq-error",
		}
	}

	return products[0], nil
}

// InsertOrderHistory insert Order History
func (s *PostgresStorage) InsertOrderHistory(ctx context.Context, orderHistory *product.OrderHistory) (*product.OrderHistory, *types.Error) {
	err := s.Storage.Insert(ctx, orderHistory)
	if err != nil {
		return nil, &types.Error{
			Path:    ".ProductPostgresStorage->Insert()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return orderHistory, nil
}

// UpdateOrderHistory update Order History
func (s *PostgresStorage) UpdateOrderHistory(ctx context.Context, orderHistory *product.OrderHistory) (*product.OrderHistory, *types.Error) {
	err := s.Storage.Update(ctx, orderHistory)
	if err != nil {
		return nil, &types.Error{
			Path:    ".ProductPostgresStorage->Update()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return orderHistory, nil
}

// DeleteOrderHistory delete a product
func (s *PostgresStorage) DeleteOrderHistory(ctx context.Context, orderHistoryID int) *types.Error {
	err := s.Storage.Delete(ctx, orderHistoryID)
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
