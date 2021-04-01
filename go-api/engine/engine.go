package engine

import (
	"github.com/gin-gonic/gin"
	"github.com/whuangz/go-example/go-api/middlewares"
	config "github.com/whuangz/go-example/go-api/config"
)

var Router *gin.Engine

func init () {
	Router = gin.New()
	Router.Use(middlewares.Cors())
}

func Connect() {
	Router.Run(config.Config.Port)
}