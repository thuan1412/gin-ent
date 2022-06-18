package service

import (
	"context"
	"gin-ent/dto"
	"gin-ent/ent"
	categoryRepo "gin-ent/ent/category"
	productRepo "gin-ent/ent/product"
	"github.com/go-redis/redis/v8"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type IProductService interface {
	GetProducts(ctx context.Context, request dto.GetProductsRequest) ([]*ent.Product, error)
}

type ProductService struct {
	Logger      *zap.Logger
	Db          *ent.Client
	RedisClient *redis.Client
}

func (p ProductService) GetProducts(ctx context.Context, request dto.GetProductsRequest) ([]*ent.Product, error) {
	query := p.Db.Product.Query()
	if request.Name != "" {
		query = query.Where(productRepo.NameContainsFold(request.Name))
	}
	if request.CategoryId != 0 {
		query = query.Where(productRepo.HasCategoryWith(categoryRepo.ID(request.CategoryId)))
	}
	return query.All(ctx)
}

func (p ProductService) GetProduct(ctx context.Context, id int) (*dto.GetProductResponse, error) {
	categoryService := CategoryService{Logger: p.Logger, Db: p.Db, RedisClient: p.RedisClient}
	var err error
	product, err := p.Db.Product.Query().Where(productRepo.ID(id)).Only(ctx)
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

	productResp.BreadCrumbs, err = categoryService.GetBreadCrumbs(ctx, category)
	productResp.BreadCrumbs = lo.Reverse[*dto.BreadCrumb](productResp.BreadCrumbs)
	return productResp, nil
}
