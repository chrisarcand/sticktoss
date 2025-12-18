package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sticktoss/backend/internal/api"
	"github.com/sticktoss/backend/internal/auth"
	"github.com/sticktoss/backend/internal/db"
	"github.com/sticktoss/backend/internal/models"
)

func main() {
	// Get database config from environment
	dbConfig := db.GetConfigFromEnv()

	// Connect to database
	database, err := db.New(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := models.Migrate(database); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database connected and migrations completed")

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	// Create router
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"}, // Vite default port
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Initialize handlers
	authHandler := api.NewAuthHandler(database)
	playerHandler := api.NewPlayerHandler(database)
	groupHandler := api.NewGroupHandler(database)

	// Public routes
	r.POST("/api/auth/signup", authHandler.Signup)
	r.POST("/api/auth/login", authHandler.Login)
	r.GET("/api/game/:shareId", groupHandler.GetGame)
	r.GET("/api/groups/:id/logo", groupHandler.GetGroupLogo)
	r.GET("/api/game/:shareId/logo", groupHandler.GetGameLogo)

	// Protected routes
	protected := r.Group("/api")
	protected.Use(auth.AuthMiddleware())
	{
		// Auth routes
		protected.GET("/auth/me", authHandler.Me)

		// Player routes
		protected.GET("/players", playerHandler.GetPlayers)
		protected.GET("/players/:id", playerHandler.GetPlayer)
		protected.POST("/players", playerHandler.CreatePlayer)
		protected.PUT("/players/:id", playerHandler.UpdatePlayer)
		protected.DELETE("/players/:id", playerHandler.DeletePlayer)

		// Group routes
		protected.GET("/groups", groupHandler.GetGroups)
		protected.GET("/groups/:id", groupHandler.GetGroup)
		protected.POST("/groups", groupHandler.CreateGroup)
		protected.PUT("/groups/:id", groupHandler.UpdateGroup)
		protected.DELETE("/groups/:id", groupHandler.DeleteGroup)

		// Group-Player routes
		protected.POST("/groups/:id/players", groupHandler.AddPlayerToGroup)
		protected.DELETE("/groups/:id/players/:player_id", groupHandler.RemovePlayerFromGroup)

		// Group logo routes
		protected.POST("/groups/:id/logo", groupHandler.UploadGroupLogo)
		protected.DELETE("/groups/:id/logo", groupHandler.DeleteGroupLogo)

		// Team generation
		protected.POST("/groups/:id/generate-teams", groupHandler.GenerateTeams)
	}

	// Serve static files from frontend build (for production)
	r.Static("/assets", "./frontend/dist/assets")
	r.StaticFile("/logo.png", "./frontend/dist/logo.png")
	r.StaticFile("/", "./frontend/dist/index.html")
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
