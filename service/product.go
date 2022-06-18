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

func (p ProductService) GetProduct(ctx context.Context, id int) (*dto.GetProductResponse, error) {
	product, err := p.Db.Product.Query().Where(product.ID(id)).Only(ctx)
	if err != nil {
		p.Logger.Warn("failed to get product", zap.Error(err))
		return nil, err
	}
	productResp := &dto.GetProductResponse{}
	productResp.ID = product.ID
	productResp.Name = product.Name
	productResp.Price = product.Price
	category, err := product.QueryCategory().Only(ctx)
	if err != nil {
		p.Logger.Warn("failed to get category", zap.Error(err))
		return nil, err
	}
	productResp.CategoryLevel = []dto.CategoryLevel{{ID: category.ID, Name: category.Name}}

	var parentCat *ent.Category
	for {
		parentCat, err = category.QueryParent().Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			p.Logger.Warn("failed to get parent category", zap.Error(err))
			return nil, err
		}
		if parentCat == nil {
			break
		}
		category = parentCat
		parentCat = nil
		productResp.CategoryLevel = append(productResp.CategoryLevel, dto.CategoryLevel{ID: category.ID, Name: category.Name})
	}
	return productResp, nil
}
