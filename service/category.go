package service

import (
	"context"
	"fmt"
	"gin-ent/dto"
	"gin-ent/ent"
	categoryRepo "gin-ent/ent/category"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"strconv"
)

type ICategoryService interface {
	GetAncestors(ctx context.Context, cat *ent.Category) ([]*ent.Category, error)
	GetBreadCrumbs(ctx context.Context, cat *ent.Category) ([]*dto.BreadCrumb, error)
}

type CategoryService struct {
	Logger      *zap.Logger
	Db          *ent.Client
	RedisClient *redis.Client
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

// DumpBreadCrumbsToRedis dump the breadCrumbs data to redis
// if catIDs is nil, dump breadCrumbs of all product
func (c CategoryService) DumpBreadCrumbsToRedis(ctx context.Context, catIDs ...int) error {
	// TODO: add option to delete key before dump
	var categories []*ent.Category
	if len(catIDs) == 0 {
		categories, _ = c.Db.Category.Query().All(ctx)
	} else {
		categories, _ = c.Db.Category.Query().Where(categoryRepo.IDIn(catIDs...)).All(ctx)
	}
	for _, category := range categories {
		breadCrumbs, err := c.getBreadCrumbsDb(ctx, category)
		if err != nil {
			return err
		}
		for idx, breadCrumb := range breadCrumbs {
			res := c.RedisClient.HSet(ctx, strconv.Itoa(category.ID), idx, breadCrumb)
			fmt.Println(res)
		}
		c.Logger.Info("Dump breadCrumbs to redis", zap.Int("category_id", category.ID), zap.Any("breadCrumbs", breadCrumbs))
	}
	return nil
}
