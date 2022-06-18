package controller

import (
	"gin-ent/helpers"
	"gin-ent/service"
	"github.com/gin-gonic/gin"
)

// DumpCategoryBreadCrumb
// @Schemes
// @Description DumpCategoryBreadCrumb
// @Tags Backdoor
// @Accept json
// @Produce json
// @Param category_id path string true "Category ID"
// @Success 200 {string} Category
// @Router /backdoor/dump-breadcrumb [POST]
func DumpCategoryBreadCrumb(ctx *gin.Context) {
	//categoryId := ctx.Param("category_id")
	//categoryIdInt, err := strconv.Atoi(categoryId)
	db := helpers.GetDbFromContext(ctx)
	redisClient := helpers.GetRedisFromContext(ctx)
	logger := helpers.GetLoggerFromContext(ctx)
	categoryService := service.CategoryService{Db: db, RedisClient: redisClient, Logger: logger}
	err := categoryService.DumpBreadCrumbsToRedis(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"status": "success"})
}
