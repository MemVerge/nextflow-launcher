package api

import "github.com/gin-gonic/gin"

type Status struct {
	Status string `json:"status" example:"ok"`
}

// Health check endpoint
// Add swagger metadata to the health check endpoint
// @Summary Health check endpoint
// @Description Check the health of the API
// @Accept  json
// @Produce json
// @Success 200 {object} Status
// @Router /health [get]
func Health(c *gin.Context) {
	s := Status{Status: "ok"}
	c.JSON(200, s)
}
