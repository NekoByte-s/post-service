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

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.POST("/posts", postHandler.CreatePost)
		v1.GET("/posts", postHandler.GetAllPosts)
		v1.GET("/posts/:id", postHandler.GetPost)
		v1.PUT("/posts/:id", postHandler.UpdatePost)
		v1.DELETE("/posts/:id", postHandler.DeletePost)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}