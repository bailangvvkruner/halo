package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/halo-dev/halo-go/internal/data"
	"github.com/halo-dev/halo-go/internal/extension"
	"github.com/halo-dev/halo-go/internal/model"
)

type SetupHandler struct {
	store *data.ExtensionStore
}

func NewSetupHandler(store *data.ExtensionStore) *SetupHandler {
	return &SetupHandler{store: store}
}

type SetupStatusResponse struct {
	Setup bool `json:"setup"`
}

type SetupRequest struct {
	Username    string `json:"username" binding:"required,min=2,max=64"`
	Password    string `json:"password" binding:"required,min=6,max=64"`
	Email       string `json:"email" binding:"required,email"`
	DisplayName string `json:"displayName"`
}

func (h *SetupHandler) GetStatus(c *gin.Context) {
	ctx := c.Request.Context()
	isInitialized, err := h.checkInitialized(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "检查初始化状态失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    SetupStatusResponse{Setup: isInitialized},
	})
}

func (h *SetupHandler) DoSetup(c *gin.Context) {
	ctx := c.Request.Context()

	isInitialized, err := h.checkInitialized(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "检查初始化状态失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	if isInitialized {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "系统已完成初始化，请勿重复操作",
			"data":    nil,
		})
		return
	}

	var req SetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	displayName := req.DisplayName
	if displayName == "" {
		displayName = req.Username
	}

	adminUser := &model.User{}
	adminUser.Metadata.Name = "admin"
	adminUser.Spec.UserName = req.Username
	adminUser.Spec.Email = req.Email
	adminUser.Spec.DisplayName = displayName
	adminUser.Spec.RawPassword = req.Password

	if _, err := h.store.Create(ctx, adminUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建管理员账号失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "系统初始化成功",
		"data":    nil,
	})
}

func (h *SetupHandler) checkInitialized(ctx context.Context) (bool, error) {
	result, err := h.store.List(ctx, &extension.ListOptions{Size: 1})
	if err != nil {
		return false, fmt.Errorf("查询用户列表失败: %w", err)
	}
	for _, ext := range result.Items {
		if _, ok := ext.(*model.User); ok {
			return true, nil
		}
	}
	return false, nil
}
