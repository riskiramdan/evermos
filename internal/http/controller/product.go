package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/riskiramdan/evermos/internal/data"
	"github.com/riskiramdan/evermos/internal/http/response"
	"github.com/riskiramdan/evermos/internal/product"
	"github.com/riskiramdan/evermos/internal/types"
)

// ProductController represents the product controller
type ProductController struct {
	productService product.ServiceInterface
	dataManager    *data.Manager
}

// ProductList product list and count
type ProductList struct {
	Data  []*product.Product `json:"data"`
	Count int                `json:"count"`
}

// ListProduct Function for listing data product
func (a *ProductController) ListProduct(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	queryValues := r.URL.Query()
	var limit = 10
	var errConversion error
	if queryValues.Get("limit") != "" {
		limit, errConversion = strconv.Atoi(queryValues.Get("limit"))
		if errConversion != nil {
			err = &types.Error{
				Path:    ".ProductController->ListProduct()",
				Message: errConversion.Error(),
				Error:   errConversion,
				Type:    "golang-error",
			}
			response.Error(w, "Bad Request", http.StatusBadRequest, *err)
			return
		}
	}

	var page = 1
	if queryValues.Get("page") != "" {
		page, errConversion = strconv.Atoi(queryValues.Get("page"))
		if errConversion != nil {
			err = &types.Error{
				Path:    ".ProductController->ListProduct()",
				Message: errConversion.Error(),
				Error:   errConversion,
				Type:    "golang-error",
			}
			response.Error(w, "Bad Request", http.StatusBadRequest, *err)
			return
		}
	}

	var search = queryValues.Get("search")

	if limit < 0 {
		limit = 10
	}
	if page < 0 {
		page = 1
	}
	productList, count, err := a.productService.ListProducts(r.Context(), &product.FindAllProductsParams{
		Limit:  limit,
		Search: search,
		Page:   page,
	})
	if err != nil {
		err.Path = ".ProductController->ListProduct()" + err.Path
		if err.Error != data.ErrNotFound {
			response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
			return
		}
	}
	if productList == nil {
		productList = []*product.Product{}
	}

	response.JSON(w, http.StatusOK, ProductList{
		Data:  productList,
		Count: count,
	})
}

// CreateProduct Function for Create data product
func (a *ProductController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	decoder := json.NewDecoder(r.Body)

	var params *product.TransactionProductParams
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		err = &types.Error{
			Path:    ".ProductController->CreateProduct()",
			Message: errDecode.Error(),
			Error:   errDecode,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}

	errTransaction := a.dataManager.RunInTransaction(r.Context(), func(ctx context.Context) error {
		_, err = a.productService.CreateProduct(ctx, params)
		if err != nil {
			return err.Error
		}
		return nil
	})
	if errTransaction != nil {
		err.Path = ".ProductController->CreateProduct()" + err.Path
		if errTransaction == product.ErrProductAlreadyExists {
			response.Error(w, product.ErrProductAlreadyExists.Error(), http.StatusUnprocessableEntity, *err)
			return
		}
		response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	response.JSON(w, http.StatusOK, "Product Created Successfully")
}

// UpdateProduct Function for update data product
func (a *ProductController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	decoder := json.NewDecoder(r.Body)

	var params *product.TransactionProductParams
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		err = &types.Error{
			Path:    ".ProductController->UpdateProduct()",
			Message: errDecode.Error(),
			Error:   errDecode,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}
	var sProductID = chi.URLParam(r, "id")
	productID, errConversion := strconv.Atoi(sProductID)
	if errConversion != nil {
		err = &types.Error{
			Path:    ".ProductController->UpdateProduct()",
			Message: errConversion.Error(),
			Error:   errConversion,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}

	var singleProduct *product.Product
	errTransaction := a.dataManager.RunInTransaction(r.Context(), func(ctx context.Context) error {
		singleProduct, err = a.productService.UpdateProduct(ctx, productID, params)
		if err != nil {
			return err.Error
		}
		return nil
	})
	if errTransaction != nil {
		err.Path = ".ProductController->UpdateProduct()" + err.Path
		if errTransaction == data.ErrAlreadyExist {
			response.Error(w, data.ErrAlreadyExist.Error(), http.StatusUnprocessableEntity, *err)
			return
		}
		response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}
	response.JSON(w, http.StatusOK, singleProduct)

}

// DeleteProduct Function for delete data product
func (a *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var err *types.Error
	var sProductID = chi.URLParam(r, "id")
	productID, errConversion := strconv.Atoi(sProductID)
	if errConversion != nil {
		err = &types.Error{
			Path:    ".ProductController->DeleteProduct()",
			Message: errConversion.Error(),
			Error:   errConversion,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}

	errTransaction := a.dataManager.RunInTransaction(r.Context(), func(ctx context.Context) error {
		err = a.productService.DeleteProduct(ctx, productID)
		if err != nil {
			return err.Error
		}
		return nil
	})
	if errTransaction != nil {
		err.Path = ".ProductController->DeleteProduct()" + err.Path
		response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}
	response.JSON(w, http.StatusNoContent, "")
}

// CreateOrder Function for delete data order
func (a *ProductController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	decoder := json.NewDecoder(r.Body)

	var params *product.TransactionOrderHistorytParams
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		err = &types.Error{
			Path:    ".ProductController->CreateOrder()",
			Message: errDecode.Error(),
			Error:   errDecode,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}

	errTransaction := a.dataManager.RunInTransaction(r.Context(), func(ctx context.Context) error {
		_, err = a.productService.CreateOrder(ctx, params)
		if err != nil {
			return err.Error
		}
		return nil
	})
	if errTransaction != nil {
		err.Path = ".ProductController->CreateOrder()" + err.Path
		response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	response.JSON(w, http.StatusOK, "Order Created Successfully")
}

// NewProductController creates a new product controller
func NewProductController(
	productService product.ServiceInterface,
	dataManager *data.Manager,
) *ProductController {
	return &ProductController{
		productService: productService,
		dataManager:    dataManager,
	}
}
