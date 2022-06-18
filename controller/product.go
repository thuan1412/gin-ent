package controller

import (
	"fmt"
	"gin-ent/dto"
	"gin-ent/ent"
	"gin-ent/ent/category"
	"gin-ent/helpers"
	"gin-ent/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// GetProducts returns products list
// @Schemes
// @Description Get products lists
// @Tags Product
// @Accept json
// @Produce json
// @Success 200 {string} Product
// @Router /products [get]
func GetProducts(ctx *gin.Context) {
	db_, _ := ctx.Get("db")
	db := db_.(*ent.Client)
	logger_, _ := ctx.Get("logger")
	logger := logger_.(*zap.Logger)
	productService := service.ProductService{Logger: logger, Db: db}

	getProductRequest := dto.GetProductsRequest{}
	if err := ctx.ShouldBind(&getProductRequest); err != nil {
		logger.Warn("error binding json", zap.Error(err))
		ctx.JSON(400, gin.H{"error": "invalid data"})
		return
	}

	products, err := productService.GetProducts(ctx, getProductRequest)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"products": products})
}

// CreateProduct creates a new product
// @Schemes
// @Description Create a new product
// @Tags Product
// @Accept json
// @Produce json
// @Param product body Product true "Product"
// @Success 200 {string} Product
// @Router /products [post]
func CreateProduct(ctx *gin.Context) {
	db_, _ := ctx.Get("db")
	logger_, _ := ctx.Get("logger")
	db := db_.(*ent.Client)
	logger := logger_.(*zap.Logger)
	var productDto dto.CreateProductRequest
	if err := ctx.ShouldBindJSON(&productDto); err != nil {
		logger.Warn("error binding json", zap.Error(err))
		ctx.JSON(400, gin.H{"error": fmt.Sprintf("invalid data: %s", err.Error())})
		return
	}

	catExist, err := db.Category.Query().Where(category.ID(productDto.CategoryId)).Exist(ctx)
	if !catExist {
		ctx.JSON(400, gin.H{"error": fmt.Sprintf("category %d not found", productDto.CategoryId)})
		return
	}
	if err != nil {
		logger.Warn("error checking category exist", zap.Error(err))
		ctx.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	product, err := db.Product.Create().SetName(productDto.Name).SetPrice(productDto.Price).SetCategoryID(productDto.CategoryId).Save(ctx)
	if err != nil {
		logger.Warn("failed to create product", zap.Error(err))
		if ent.IsConstraintError(err) {
			ctx.JSON(400, gin.H{"error": fmt.Sprintf("invalid data: %s", err.Error())})
			return
		}
		ctx.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	ctx.JSON(200, gin.H{"product": product})
}

// GetProduct returns a product
// @Schemes
// @Description Get a product
// @Tags Product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {string} Product
// @Router /products/{id} [get]
func GetProduct(ctx *gin.Context) {
	db_, _ := ctx.Get("db")
	db := db_.(*ent.Client)
	logger_, _ := ctx.Get("logger")
	logger := logger_.(*zap.Logger)
	redisClient := helpers.GetRedisFromContext(ctx)
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Warn("error converting id", zap.Error(err))
		ctx.JSON(400, gin.H{"error": fmt.Sprintf("invalid data: %s", err.Error())})
		return
	}

	// TODO: create generic function to get service from context
	productService := service.ProductService{Logger: logger, Db: db, RedisClient: redisClient}
	product, err := productService.GetProduct(ctx, id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"product": product})
}
