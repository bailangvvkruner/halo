package handler

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

type ActuatorHandler struct{}

func NewActuatorHandler() *ActuatorHandler {
	return &ActuatorHandler{}
}

func (h *ActuatorHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
		"components": gin.H{
			"db": gin.H{"status": "UP"},
			"api": gin.H{"status": "UP"},
		},
	})
}

func (h *ActuatorHandler) Info(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"app": gin.H{
			"name":        "halo-go",
			"version":     "0.1.0",
			"description": "Halo CMS Blog System - Go Clone",
		},
		"runtime": gin.H{
			"name":    runtime.GOOS,
			"version": runtime.Version(),
			"arch":    runtime.GOARCH,
		},
	})
}
