package controller

import (
	"fmt"
	"gin-ent/ent"
	"github.com/gin-gonic/gin"
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
	db_, exist := ctx.Get("db")
	if !exist {
		ctx.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	db := db_.(*ent.Client)
	products, err := db.Product.Query().All(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "internal server error"})
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
	db_, exist := ctx.Get("db")
	if !exist {
		ctx.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	db := db_.(*ent.Client)
	var product ent.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		fmt.Println("failed to bind json:", err)
		ctx.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	p, err := db.Product.Create().SetName(product.Name).SetPrice(product.Price).Save(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	ctx.JSON(200, gin.H{"product": p})
}
