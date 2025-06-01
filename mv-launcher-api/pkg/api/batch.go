package api

import (
	"github.com/MemVerge/nf-launcher/pkg/types"
	"github.com/aws/aws-sdk-go-v2/service/batch"
	"github.com/gin-gonic/gin"
)

// @Summary List all queues
// @Description Returns a JSON blob with a list of all queues
// @Accept  json
// @Produce json
// @Success 200 {object} BatchQueues
// @Router /batch/queues [get]
func (a API) ListQueues(c *gin.Context) {
	out, err := a.batchClient.DescribeJobQueues(c.Request.Context(), &batch.DescribeJobQueuesInput{})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	qs := types.Queues{}
	for _, q := range out.JobQueues {
		qs = append(qs, types.Queue{
			Name: *q.JobQueueName,
			ARN:  *q.JobQueueArn,
		})
	}
	c.JSON(200, gin.H{"queues": qs})
}
