package engine

import (
	"github.com/gin-gonic/gin"
	config "github.com/whuangz/go-example/go-api/config"
	"github.com/whuangz/go-example/go-api/middlewares"
)

var Router *gin.Engine

func init() {
	Router = gin.New()
	Router.Use(middlewares.Cors())
}

func Connect() {
	Router.Run(config.PORT)
}
