package api

import (
	"context"
	"log"

	"github.com/ChristianKniep/mv-launcher-api/pkg/sss"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gin-gonic/gin"
)

// @Summary List all available pipelines
// @Description Returns a JSON blob with a list of all jobs
// @Accept  json
// @Produce json
// @Success 200 {object} types.Pipelines
// @Router /pipeline [get]
func (a API) ListPipelines(c *gin.Context) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(AwsRegion))
	if err != nil {
		log.Fatalf("failed to create session, %v", err)
	}

	pipelines, err := sss.GetPipelines(cfg, PipelineBucket)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, pipelines)
}
