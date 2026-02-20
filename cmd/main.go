package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/user/debt-optimization-engine/config"
	"github.com/user/debt-optimization-engine/internal/handlers"
	"github.com/user/debt-optimization-engine/internal/repositories"
	"github.com/user/debt-optimization-engine/internal/services"
)

func main() {
	// 1. Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Could not load config:", err)
	}

	// 2. Database Connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		log.Fatalf("Unable to connect to database at %s: %v", cfg.DBURL, err)
	}
	defer pool.Close()

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	// 3. Initialize Layers
	repo := repositories.NewPostgresRepo(pool)
	settlementSvc := services.NewSettlementService(repo)
	h := handlers.NewHandler(repo, settlementSvc)

	// 4. Setup Router
	r := gin.New() // Use New() to manually add middleware

	// Global Middleware
	r.Use(gin.Logger())
	
	// Structured Panic Recovery
	r.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC RECOVERED: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
					"status": 500,
				})
			}
		}()
		c.Next()
	})

	// Routes
	api := r.Group("")
	{
		api.POST("/users", h.CreateUser)
		api.POST("/groups", h.CreateGroup)
		api.POST("/groups/:id/members", h.AddMember)
		api.POST("/groups/:id/expenses", h.CreateExpense)
		api.GET("/groups/:id/balances", h.GetBalances)
		api.GET("/groups/:id/settlement", h.GetSettlement)
		api.GET("/groups/:id/settlement/compare", h.CompareStrategies)
	}

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		dbStatus := "connected"
		if err := pool.Ping(c.Request.Context()); err != nil {
			dbStatus = "disconnected"
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"database": dbStatus,
			"time":     time.Now().Format(time.RFC3339),
		})
	})

	log.Printf("Server starting on port %s", cfg.Port)
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
