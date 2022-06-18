package service

import (
	"context"
	"gin-ent/dto"
	"gin-ent/ent"
	"gin-ent/ent/category"
	"gin-ent/ent/product"
	"go.uber.org/zap"
)

type IProductService interface {
	GetProducts(ctx context.Context, request dto.GetProductsRequest) ([]*ent.Product, error)
}

type ProductService struct {
	Logger *zap.Logger
	Db     *ent.Client
}

func (p ProductService) GetProducts(ctx context.Context, request dto.GetProductsRequest) ([]*ent.Product, error) {
	query := p.Db.Product.Query()
	if request.Name != "" {
		query = query.Where(product.NameContainsFold(request.Name))
	}
	if request.CategoryId != 0 {
		query = query.Where(product.HasCategoryWith(category.ID(request.CategoryId)))
	}
	return query.All(ctx)
}
