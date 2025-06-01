package api

import (
	"github.com/MemVerge/nf-launcher/pkg/services"
	"github.com/gin-gonic/gin"
)

// @Summary List all available pipelines
// @Description Returns a JSON blob with a list of all jobs
// @Accept  json
// @Produce json
// @Success 200 {object} types.Pipelines
// @Router /pipeline [get]
func (a API) ListPipelines(c *gin.Context) {
	pipelines, err := services.GetPipelines(a.s3Client, a.config.PipelineBucket)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, pipelines)
}
