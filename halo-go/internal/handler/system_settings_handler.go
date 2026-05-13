package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/halo-dev/halo-go/internal/data"
	"github.com/halo-dev/halo-go/internal/extension"
	"github.com/halo-dev/halo-go/internal/model"
)

type SystemSettingsHandler struct {
	store *data.ExtensionStore
}

func NewSystemSettingsHandler(store *data.ExtensionStore) *SystemSettingsHandler {
	return &SystemSettingsHandler{store: store}
}

func (h *SystemSettingsHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	result, err := h.store.List(ctx, &extension.ListOptions{
		Size: 0,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取系统设置失败: " + err.Error(),
		})
		return
	}
	for _, item := range result.Items {
		if settings, ok := item.(*model.SystemSettings); ok {
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "ok",
				"data":    settings,
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    model.SystemSettings{},
	})
}

func (h *SystemSettingsHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var req model.SystemSpec
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	result, err := h.store.List(ctx, &extension.ListOptions{Size: 0})
	if err != nil || len(result.Items) == 0 {
		settings := &model.SystemSettings{}
		settings.Metadata.Name = "system-settings"
		settings.Spec = req
		h.store.Create(ctx, settings)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "系统设置保存成功",
			"data":    settings,
		})
		return
	}

	for _, item := range result.Items {
		if settings, ok := item.(*model.SystemSettings); ok {
			settings.Spec = req
			updated, err := h.store.Update(ctx, settings)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "更新系统设置失败: " + err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "系统设置更新成功",
				"data":    updated,
			})
			return
		}
	}
}
