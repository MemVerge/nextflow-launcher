package main

import (
	"embed"
	"log"

	"github.com/ChristianKniep/mv-launcher-api/pkg/api"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:embed vue/dist/*
var embeddedFiles embed.FS

// @contact.name   GTM MemVerge Inc.

func main() {
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Serve static files
	router.Static("/static", "./vue/dist")
	router.StaticFile("/", "./vue/dist/index.html")
	router.StaticFile("/favicon.ico", "./vue/dist/favicon.ico")

	// Create API instance
	apiInstance := api.NewAPI()

	// API routes
	v1 := router.Group("/v1")
	{
		// Health check
		v1.GET("/health", api.Health)

		// Bucket routes
		buckets := v1.Group("/buckets")
		{
			buckets.GET("", apiInstance.ListBuckets)
		}

		// Pipeline routes
		pipelines := v1.Group("/pipelines")
		{
			pipelines.GET("", apiInstance.ListPipelines)
		}

		// Job routes
		jobs := v1.Group("/jobs")
		{
			jobs.GET("", apiInstance.ListJobs)
			jobs.POST("", apiInstance.CreateJob)
			jobs.GET("/:id/logs", apiInstance.GetJobLogs)
			jobs.GET("/:id/log-url", apiInstance.GetJobLogPresignedURL)
		}

		// Batch routes
		batch := v1.Group("/batch")
		{
			batch.GET("/queues", apiInstance.ListQueues)
			batch.POST("/setup", apiInstance.SetupAWSBatch)
		}
	}

	// use ginSwagger middleware to serve the API docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
