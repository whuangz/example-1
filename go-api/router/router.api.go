package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/whuangz/go-example/go-api/config"
	"github.com/whuangz/go-example/go-api/helpers/db"
	"github.com/whuangz/go-example/go-api/middleware"
)

var (
	router   *gin.Engine
	database *sqlx.DB
)

func init() {

	router = gin.New()
	router.Use(middleware.Cors())

	database = db.Configure("")
	err := database.Ping()
	if err != nil {
		log.Fatal("Database Open Connection: ", err)
	}

}

func Routing() {

	defer func() {
		err := database.Close()
		if err != nil {
			log.Fatal("Database Close Connection: ", err)
		}
	}()

	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "I m good",
		})
	})

	blogRoutes()

	router.Run(config.PORT)
}
