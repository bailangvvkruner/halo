package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/halo-dev/halo-go/internal/model"
	"github.com/halo-dev/halo-go/internal/service"
)

type SinglePageHandler struct {
	svc service.SinglePageService
}

func NewSinglePageHandler(svc service.SinglePageService) *SinglePageHandler {
	return &SinglePageHandler{svc: svc}
}

func (h *SinglePageHandler) Create(c *gin.Context) {
	var sp model.SinglePage
	if err := c.ShouldBindJSON(&sp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误: " + err.Error(), "data": nil})
		return
	}
	created, err := h.svc.Create(c.Request.Context(), &sp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建页面失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": created})
}

func (h *SinglePageHandler) Get(c *gin.Context) {
	name := c.Param("name")
	sp, err := h.svc.Get(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "页面不存在: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": sp})
}

func (h *SinglePageHandler) Update(c *gin.Context) {
	name := c.Param("name")
	var sp model.SinglePage
	if err := c.ShouldBindJSON(&sp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误: " + err.Error(), "data": nil})
		return
	}
	sp.Metadata.Name = name
	updated, err := h.svc.Update(c.Request.Context(), &sp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新页面失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": updated})
}

func (h *SinglePageHandler) Delete(c *gin.Context) {
	name := c.Param("name")
	if err := h.svc.Delete(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除页面失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": nil})
}

func (h *SinglePageHandler) List(c *gin.Context) {
	opts := parseListOptions(c)
	result, err := h.svc.List(c.Request.Context(), &opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询页面列表失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": result})
}

func (h *SinglePageHandler) Publish(c *gin.Context) {
	name := c.Param("name")
	if err := h.svc.Publish(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "发布页面失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": nil})
}

func (h *SinglePageHandler) Unpublish(c *gin.Context) {
	name := c.Param("name")
	if err := h.svc.Unpublish(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "取消发布失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": nil})
}

func (h *SinglePageHandler) Trash(c *gin.Context) {
	name := c.Param("name")
	if err := h.svc.Trash(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "移入回收站失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": nil})
}

func (h *SinglePageHandler) Restore(c *gin.Context) {
	name := c.Param("name")
	if err := h.svc.Restore(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "恢复页面失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": nil})
}
