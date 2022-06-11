package controller

import (
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
