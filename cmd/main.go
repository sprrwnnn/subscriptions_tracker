package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/sprrwnnn/subscription-tracker/docs"
	"github.com/sprrwnnn/subscription-tracker/internal/config"
	"github.com/sprrwnnn/subscription-tracker/internal/handler"
	"github.com/sprrwnnn/subscription-tracker/internal/middleware"
	"github.com/sprrwnnn/subscription-tracker/internal/repository"
	"github.com/sprrwnnn/subscription-tracker/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Subscription Service API
// @version 1.0
// @description Service for managing subscriptions
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Setup logger
	logrus.SetFormatter(&logrus.JSONFormatter{})
	level, _ := logrus.ParseLevel(cfg.LogLevel)
	logrus.SetLevel(level)

	// Connect to database
	db, err := gorm.Open(postgres.Open(cfg.GetDBDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		logrus.Fatal("Failed to connect to database:", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		logrus.Fatal("Failed to get database handle:", err)
	}

	// Setup dependencies
	subRepo := repository.NewSubscriptionRepository(db, logrus.StandardLogger())
	subService := service.NewSubscriptionService(subRepo, logrus.StandardLogger())
	subHandler := handler.NewSubscriptionHandler(subService, logrus.StandardLogger())

	// Setup router
	router := gin.Default()
	router.Use(middleware.LoggingMiddleware(logrus.StandardLogger()))
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	router.GET("/ready", func(c *gin.Context) {
		if err := sqlDB.PingContext(c.Request.Context()); err != nil {
			logrus.WithError(err).Error("Readiness check failed")
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "error",
				"error":  "database is not ready",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	{
		api.POST("/subscriptions", subHandler.Create)
		api.GET("/subscriptions/:id", subHandler.GetByID)
		api.PUT("/subscriptions/:id", subHandler.Update)
		api.DELETE("/subscriptions/:id", subHandler.Delete)
		api.GET("/subscriptions", subHandler.List)
		api.POST("/calculate-cost", subHandler.CalculateCost)
	}

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	logrus.WithField("addr", addr).Info("Starting server")
	if err := router.Run(addr); err != nil {
		logrus.Fatal("Failed to start server:", err)
	}
}
