package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"stars/model"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	envPtr := os.Getenv("RUN_AS")

	ambiente := ""
	switch envPtr {
	case "production":
		if err := godotenv.Load(".env.production"); err != nil {
			log.Print("No .env file found")
		}
		ambiente = envPtr

	case "development":
		if err := godotenv.Load(".env.development"); err != nil {
			log.Print("No .env file found")
		}
		ambiente = envPtr

	default:
		if err := godotenv.Load(".env"); err != nil {
			log.Print("No .env file found")
		}
		ambiente = "local"
	}
	version := os.Getenv("VERSION")
	log.Print("Running Stars ", version)
	log.Print("Environment: ", ambiente)

}

func main() {
	r := gin.Default()

	// Logging to a file, by default gin.DefaultWriter = os.Stdout
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	// test route
	r.GET("/ping", func(c *gin.Context) {
		start := time.Now()
		t := time.Now()
		elapsed := t.Unix() - start.Unix()
		c.JSON(http.StatusOK, gin.H{
			"message":    "pong",
			"time start": t,
			"duration":   elapsed,
		})
	})

	r.GET("/system", func(c *gin.Context) {
		main_stellar_object := model.Sunme()

		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"stellar": main_stellar_object,
		})
	})

	port := os.Getenv("PORT")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Initializing the server in a goroutine so that it won't block the graceful shutdown handling below
	go func() {
		log.Printf("listen on: %s\n", srv.Addr)
		if err := srv.ListenAndServeTLS("./server.crt", "./server.key"); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server quiting")
}
