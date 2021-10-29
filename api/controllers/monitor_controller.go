package controllers

import (
	"github.com/free-health/health24-gateway/api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) GetAllMonitors(c *gin.Context) {
	monitors, err := models.Monitor{}.GetAll(s.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"monitors": monitors,
	})
}
