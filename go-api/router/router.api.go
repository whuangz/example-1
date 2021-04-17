package router

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	accountRoutes()
	blogRoutes()

	srv := &http.Server{
		Addr:    config.PORT,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialize server: %v\n", err)
		}
	}()

	log.Printf("Listening on port %v\n", srv.Addr)

	// Wait for kill signal of channel
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This blocks until a signal is passed into the quit channel
	<-quit

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server
	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}
