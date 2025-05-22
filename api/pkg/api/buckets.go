package api

import (
	"context"
	"log"

	"github.com/ChristianKniep/mv-launcher-api/pkg/sss"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gin-gonic/gin"
)

func (a API) ListBuckets(c *gin.Context) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(AwsRegion))
	if err != nil {
		log.Fatalf("failed to create session, %v", err)
	}

	buckets, err := sss.ListBuckets(cfg)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, buckets)
}
