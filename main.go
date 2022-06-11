package main

import (
	"context"
	"entgo.io/ent/dialect"
	"fmt"
	docs "gin-ent/docs"
	"gin-ent/ent"
	"gin-ent/helpers"
	"gin-ent/route"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

func InjectDbClient() gin.HandlerFunc {
	return func(c *gin.Context) {
		entClient, err := helpers.GetDb()
		if err != nil {
			fmt.Println("failed to open connection to database:", err)
			c.JSON(500, gin.H{"error": "internal server error"})
		}
		c.Set("db", entClient)
	}
}

func main() {
	// start gin server
	r := gin.Default()
	r.Use(InjectDbClient())

	entClient, err := ent.Open(dialect.Postgres,
		"host=localhost port=5432 user=postgres dbname=gin-ent password=1 sslmode=disable")
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
