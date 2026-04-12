package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
