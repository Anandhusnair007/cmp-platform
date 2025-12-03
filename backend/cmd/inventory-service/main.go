package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/cmp-platform/backend/internal/config"
	"github.com/cmp-platform/backend/internal/db"
	"github.com/cmp-platform/backend/internal/handlers"
	"github.com/cmp-platform/backend/internal/middleware"
)

var (
	port = flag.Int("port", 8081, "Server port")
)

func main() {
	flag.Parse()

	cfg := config.Load()
	if *port != 8081 {
		cfg.ServerPort = *port
	}

	// Initialize database
	database, err := db.New()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create handler
	handler := handlers.NewInventoryHandler(database)

	// Setup router
	router := gin.Default()
	
	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	
	// Health check (no auth required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "inventory-service"})
	})

	// Metrics endpoint (no auth required)
	router.GET("/metrics", func(c *gin.Context) {
		c.String(http.StatusOK, "# Inventory service metrics\n")
	})

	// API routes with authentication
	api := router.Group("/api/v1")
	api.Use(middleware.RequireAuth())
	{
		api.GET("/inventory", handler.GetInventory)
		api.GET("/inventory/expiring", handler.GetExpiringCertificates)
	}

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,
	}

	go func() {
		log.Printf("Inventory service starting on port %d", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down inventory service...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Inventory service exited")
}
