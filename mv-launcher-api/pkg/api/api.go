package api

import (
	"github.com/MemVerge/nf-launcher/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/batch"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

// API represents the API server
type API struct {
	config      *config.Config
	batchClient *batch.Client
	s3Client    *s3.Client
}

// NewAPI creates a new API instance
func NewAPI(cfg *config.Config, batchClient *batch.Client, s3Client *s3.Client) *API {
	return &API{
		config:      cfg,
		batchClient: batchClient,
		s3Client:    s3Client,
	}
}

// RegisterRoutes registers all API routes
func (a *API) RegisterRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", a.Health)

	// API routes
	v1 := router.Group("/v1")
	{
		// Bucket routes
		buckets := v1.Group("/buckets")
		{
			buckets.GET("", a.ListBuckets)
		}

		// Pipeline routes
		pipelines := v1.Group("/pipelines")
		{
			pipelines.GET("", a.ListPipelines)
		}

		// Job routes
		jobs := v1.Group("/jobs")
		{
			jobs.GET("", a.ListJobs)
			jobs.POST("", a.CreateJob)
			jobs.GET("/:id/logs", a.GetJobLogs)
			jobs.GET("/:id/log-url", a.GetJobLogPresignedURL)
		}

		// Batch routes
		batch := v1.Group("/batch")
		{
			batch.GET("/queues", a.ListQueues)
		}
	}
}

// Health handles health check requests
func (a *API) Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
