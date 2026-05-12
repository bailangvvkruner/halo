package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/halo-dev/halo-go/internal/model"
	"github.com/halo-dev/halo-go/internal/service"
)

type MenuItemHandler struct {
	svc service.MenuItemService
}

func NewMenuItemHandler(svc service.MenuItemService) *MenuItemHandler {
	return &MenuItemHandler{svc: svc}
}

func (h *MenuItemHandler) Create(c *gin.Context) {
	var item model.MenuItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误: " + err.Error(), "data": nil})
		return
	}
	created, err := h.svc.Create(c.Request.Context(), &item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建菜单项失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": created})
}

func (h *MenuItemHandler) Get(c *gin.Context) {
	name := c.Param("name")
	item, err := h.svc.Get(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "菜单项不存在: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": item})
}

func (h *MenuItemHandler) Update(c *gin.Context) {
	name := c.Param("name")
	var item model.MenuItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误: " + err.Error(), "data": nil})
		return
	}
	item.Metadata.Name = name
	updated, err := h.svc.Update(c.Request.Context(), &item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新菜单项失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": updated})
}

func (h *MenuItemHandler) Delete(c *gin.Context) {
	name := c.Param("name")
	if err := h.svc.Delete(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除菜单项失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": nil})
}

func (h *MenuItemHandler) List(c *gin.Context) {
	opts := parseListOptions(c)
	result, err := h.svc.List(c.Request.Context(), &opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询菜单项列表失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": result})
}
