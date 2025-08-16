package main

import (
	"log"
	"os"

	"postService/internal/database"
	"postService/internal/handlers"
	"postService/internal/repository"
	"postService/internal/service"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "postService/docs"
)

// @title Post Service API
// @version 1.0
// @description A simple post service with CRUD operations
// @host localhost:8080
// @BasePath /api/v1
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize database
	dbConfig := database.NewConfig()
	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize repository with PostgreSQL
	repo := repository.NewPostgresPostRepository(db.DB)
	postService := service.NewPostService(repo)
	postHandler := handlers.NewPostHandler(postService)

	// Initialize health service and handler
	healthService := service.NewHealthService(db, "1.0.0")
	healthHandler := handlers.NewHealthHandler(healthService)

	router := gin.Default()

	// API endpoints
	v1 := router.Group("/api/v1")
	{
		// Post endpoints
		v1.POST("/posts", postHandler.CreatePost)
		v1.GET("/posts", postHandler.GetAllPosts)
		v1.GET("/posts/:id", postHandler.GetPost)
		v1.PUT("/posts/:id", postHandler.UpdatePost)
		v1.DELETE("/posts/:id", postHandler.DeletePost)
		
		// Health endpoints
		v1.GET("/health", healthHandler.GetHealth)           // GET /api/v1/health
		v1.GET("/health/live", healthHandler.GetLiveness)    // GET /api/v1/health/live
		v1.GET("/health/ready", healthHandler.GetReadiness)  // GET /api/v1/health/ready
		v1.GET("/health/ping", healthHandler.GetHealthSimple) // GET /api/v1/health/ping
		v1.GET("/health/component/:component", healthHandler.GetComponentHealth) // GET /api/v1/health/component/{name}
	}

	// Keep infrastructure health endpoints for Kubernetes probes (no versioning)
	router.GET("/health/live", healthHandler.GetLiveness)
	router.GET("/health/ready", healthHandler.GetReadiness)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}