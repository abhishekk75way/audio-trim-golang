package main

import (
	"authentication/backend/internal/config"
	"authentication/backend/internal/handlers"
	"authentication/backend/internal/middleware"
	"authentication/backend/internal/models"
	"authentication/backend/internal/queue"
	"authentication/backend/internal/repositories"
	"authentication/backend/internal/routes"
	"authentication/backend/internal/services"
	"strings"

	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "authentication/backend/docs"

	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

// @title API Docs
// @version 1.0
// @description API documentation for Authentication & Authorization system
// @termsOfService http://swagger.io/terms/
// @contact.name Abhishek Kushwaha
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @schemes http https
func main() {
	originsEnv := os.Getenv("CORS_ORIGINS")
	var allowedOrigins []string

	if originsEnv == "" {
		allowedOrigins = []string{"http://localhost:5173"}
	} else {
		for _, o := range strings.Split(originsEnv, ",") {
			if v := strings.TrimSpace(o); v != "" {
				allowedOrigins = append(allowedOrigins, v)
			}
		}
	}

	str := os.Getenv("POSTGRES_STR")
	if str == "" {
		str = "host=localhost user=postgres password=postgres dbname=authdb port=5432 sslmode=disable"
	}

	if err := config.Connect(str); err != nil {
		log.Fatal("DB connect failed:", err)
	}

	if err := config.DB.AutoMigrate(
		&models.User{},
		&models.Job{},
		&models.JobFile{},
	); err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	config.ConnectRedis()

	userRepo := repositories.NewUserRepo(config.DB)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	jobRepo := repositories.NewJobRepo(config.DB)
	jobFileRepo := repositories.NewJobFileRepo(config.DB)

	jobQueue := queue.NewQueue()
	jobQueue.Start(jobRepo, jobFileRepo)

	jobHandler := &handlers.Handler{
		Jobs:  jobRepo,
		Files: jobFileRepo,
		Queue: jobQueue,
	}

	r := gin.New()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.Use(
		gin.Recovery(),
		middleware.ErrorHandler(),
		middleware.JWTOptional(),
	)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.Setup(r, authHandler, jobHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
