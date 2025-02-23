package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hussainr95/link-in-bio-service/config"
	httphandlers "github.com/hussainr95/link-in-bio-service/internal/delivery/http"
	"github.com/hussainr95/link-in-bio-service/internal/repository"
	"github.com/hussainr95/link-in-bio-service/internal/usecase"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "github.com/hussainr95/link-in-bio-service/docs"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// @title Link in Bio API
// @version 1.0
// @description API for managing Link in Bio entries.
// @host localhost:8080
// @BasePath /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description "Enter your bearer token in the format: Bearer test"

func main() {
	// 1. Load configuration.
	cfg := config.NewConfig()

	// 2. Connect to MongoDB.
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer client.Disconnect(context.Background())

	// Ping DB to verify the connection.
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal("Could not ping DB:", err)
	}

	db := client.Database(cfg.MongoDBName)

	// 3. Setup repositories.
	linkRepo := repository.NewMongoLinkRepository(db)
	visitRepo := repository.NewMongoVisitRepository(db)

	// 4. Setup usecase with both repositories.
	linkUsecase := usecase.NewLinkUsecase(linkRepo, visitRepo)

	// 5. Setup Gin router.
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Bonus step. Register Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 6. Apply authentication middleware globally.
	// All endpoints (except maybe /visit) can be protected. Adjust as needed.
	router.Use(httphandlers.AuthMiddleware())

	// 7. Register link routes.
	linkHandler := httphandlers.NewLinkHandler(linkUsecase)
	linkHandler.RegisterAPIRoutes(router)

	// 8. Start background cleanup goroutine.
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if err := linkUsecase.CleanupExpiredLinks(context.Background()); err != nil {
				log.Println("Error cleaning expired links:", err)
			} else {
				log.Println("Expired links cleaned successfully.")
			}
		}
	}()

	// 9. Start the server.
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Println("Server running on", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
