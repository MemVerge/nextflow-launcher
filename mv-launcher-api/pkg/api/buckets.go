package api

import (
	"github.com/MemVerge/nf-launcher/pkg/services"
	"github.com/gin-gonic/gin"
)

func (a API) ListBuckets(c *gin.Context) {
	buckets, err := services.ListBuckets(a.s3Client)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, buckets)
}
