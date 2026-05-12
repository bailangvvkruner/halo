package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/halo-dev/halo-go/internal/service"
)

type StatsHandler struct {
	svc service.StatsService
}

func NewStatsHandler(svc service.StatsService) *StatsHandler {
	return &StatsHandler{svc: svc}
}

func (h *StatsHandler) IncrVisit(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		name = c.Param("name")
	}
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 name 参数", "data": nil})
		return
	}
	if err := h.svc.IncrVisit(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "统计访问量失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": nil})
}

func (h *StatsHandler) GetVisitCount(c *gin.Context) {
	name := c.Param("name")
	count, err := h.svc.GetVisitCount(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取访问量失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": gin.H{"count": count}})
}

func (h *StatsHandler) GetAllStats(c *gin.Context) {
	stats, err := h.svc.GetAllStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取统计数据失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": stats})
}
