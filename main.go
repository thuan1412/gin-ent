package main

import (
	"context"
	docs "gin-ent/docs"
	"gin-ent/ent"
	"gin-ent/helpers"
	"gin-ent/route"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"log"
)

var logger *zap.Logger
var entClient *ent.Client
var redisClient *redis.Client
var err error

func init() {
	logger = zap.NewExample()
	entClient, err = helpers.GetDb()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       5,
	})
	if err != nil {
		panic(err)
	}
}

func InjectDbClient() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", entClient)
	}
}
func InjectLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("logger", logger)
		c.Set("redis", redisClient)
	}
}

func main() {
	// start gin server
	r := gin.Default()
	r.Use(InjectDbClient())
	r.Use(InjectLogger())

	entClient, err := helpers.GetDb()
	if err != nil {
		log.Fatalln("failed to open connection to database:", err)
	}

	// schema auto-migration
	defer entClient.Close()
	if err := entClient.Schema.Create(context.Background()); err != nil {
		log.Fatalln("failed to create schema:", err)
	}

	err = r.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}

	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	route.SetUpRoute(r)

	err = r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
