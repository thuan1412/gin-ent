package helpers

import (
	"gin-ent/ent"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func GetLoggerFromContext(ctx *gin.Context) *zap.Logger {
	logger_, ok := ctx.Get("logger")
	if !ok || logger_ == nil {
		panic("logger not found in context")
	}
	logger := logger_.(*zap.Logger)
	return logger
}

func GetDbFromContext(ctx *gin.Context) *ent.Client {
	db_, ok := ctx.Get("db")
	if !ok || db_ == nil {
		panic("redis not found in context")
	}
	db := db_.(*ent.Client)
	return db
}

func GetRedisFromContext(ctx *gin.Context) *redis.Client {
	redis_, ok := ctx.Get("redis")
	if !ok || redis_ == nil {
		panic("redis not found in context")
	}
	redisClient := redis_.(*redis.Client)
	return redisClient
}
