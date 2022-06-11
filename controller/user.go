package controller

import "github.com/gin-gonic/gin"

// GetUsers returns users list
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /users [ping get]
func GetUsers(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"user": "a"})
}
