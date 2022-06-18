package service

import (
	"context"
	"gin-ent/dto"
	"gin-ent/ent"
	"go.uber.org/zap"
)

type ICategoryService interface {
	GetAncestors(ctx context.Context, cat *ent.Category) ([]*ent.Category, error)
	GetBreadCrumbs(ctx context.Context, cat *ent.Category) ([]*dto.BreadCrumb, error)
}

type CategoryService struct {
	Logger *zap.Logger
	Db     *ent.Client
}

func (c CategoryService) GetAncestors(ctx context.Context, category *ent.Category) ([]*ent.Category, error) {
	var parentCat *ent.Category
	var err error
	ancestors := []*ent.Category{
		category,
	}
	for {
		parentCat, err = category.QueryParent().Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			c.Logger.Warn("failed to get parent category", zap.Error(err))
			return nil, err
		}
		if parentCat == nil {
			break
		}
		category = parentCat
		parentCat = nil
		ancestors = append(ancestors, category)
	}
	return ancestors, nil
}

func (c CategoryService) GetBreadCrumbs(ctx context.Context, category *ent.Category) ([]*dto.BreadCrumb, error) {
	return c.getBreadCrumbsDb(ctx, category)
}

// getAncestorsDb get the ancestors of the category from database
func (c CategoryService) getBreadCrumbsDb(ctx context.Context, category *ent.Category) ([]*dto.BreadCrumb, error) {
	ancestors, err := c.GetAncestors(ctx, category)
	if err != nil {
		return nil, err
	}
	breadCrumbs := make([]*dto.BreadCrumb, len(ancestors))
	for i, ancestor := range ancestors {
		breadCrumbs[i] = &dto.BreadCrumb{
			ID:   ancestor.ID,
			Name: ancestor.Name,
		}
	}
	return breadCrumbs, nil
}

// getAncestorsRedis get the ancestors of the category from database
func (c CategoryService) getBredCrumbRedis(ctx context.Context, category *ent.Category) ([]*dto.BreadCrumb, error) {
	return nil, nil
}
