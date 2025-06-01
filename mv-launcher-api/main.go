package main

import (
	"context"
	"fmt"
	"log"

	"github.com/MemVerge/nf-launcher/pkg/api"
	configlocal "github.com/MemVerge/nf-launcher/pkg/config"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/batch"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := configlocal.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Load AWS config
	awsCfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(cfg.AWSRegion),
	)
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	// Initialize AWS clients
	batchClient := batch.NewFromConfig(awsCfg)
	s3Client := s3.NewFromConfig(awsCfg)

	// Initialize API
	apiInstance := api.NewAPI(cfg, batchClient, s3Client)

	// Create router
	router := gin.Default()

	// Configure CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", cfg.CORSAllowedOrigins[0])
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Register routes
	apiInstance.RegisterRoutes(router)

	// Start server
	log.Printf("Starting server on port %d", cfg.Port)
	if err := router.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
