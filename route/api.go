package route

import (
	"gin-ent/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Ping pong
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /ping [get]
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}

// SetUpRoute sets up the routes for the application.
func SetUpRoute(router *gin.Engine) {
	router.GET("/users", controller.GetUsers)
	router.GET("/ping", ping)
}
