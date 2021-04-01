package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/whuangz/go-example/go-api/engine"
)

func init() {

	engine.Router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "I m good",
		})
	})

	blogsGroup := engine.Router.Group("blogs")
	{
		blogsGroup.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

	}
}
