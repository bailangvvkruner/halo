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
	Username   string `json:"username" binding:"required,min=4,max=63"`
	Password   string `json:"password" binding:"required,min=5,max=257"`
	Email      string `json:"email" binding:"required,email"`
	SiteTitle  string `json:"siteTitle" binding:"required,max=80"`
	Language   string `json:"language" binding:"omitempty,oneof=zh-CN zh-TW en es"`
	ExternalUrl string `json:"externalUrl" binding:"required,url"`
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

	adminUser := &model.User{}
	adminUser.Metadata.Name = "admin"
	adminUser.Spec.UserName = req.Username
	adminUser.Spec.Email = req.Email
	adminUser.Spec.DisplayName = req.Username
	adminUser.Spec.RawPassword = req.Password

	if _, err := h.store.Create(ctx, adminUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建管理员账号失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	h.initializeSystemSettings(ctx, &req)
	h.initializePresetData(ctx, req.Username)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "系统初始化成功",
		"data": gin.H{
			"redirectTo": "/console",
		},
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

func (h *SetupHandler) initializeSystemSettings(ctx context.Context, req *SetupRequest) {
	settings := &model.SystemSettings{}
	settings.Metadata.Name = "system"
	settings.Spec.Basic.Title = req.SiteTitle
	settings.Spec.Basic.Language = req.Language
	if settings.Spec.Basic.Language == "" {
		settings.Spec.Basic.Language = "zh-CN"
	}
	settings.Spec.Basic.ExternalUrl = req.ExternalUrl
	h.store.Create(ctx, settings)
}

func (h *SetupHandler) initializePresetData(ctx context.Context, username string) {
	defaultCategory := &model.Category{}
	defaultCategory.Metadata.Name = "76514a40-6ef1-4ed9-b58a-e26945bde3ca"
	defaultCategory.Spec.DisplayName = "默认分类"
	defaultCategory.Spec.Slug = "default"
	defaultCategory.Spec.Description = "这是你的默认分类，如不需要，删除即可。"
	h.store.Create(ctx, defaultCategory)

	haloTag := &model.Tag{}
	haloTag.Metadata.Name = "c33ceabb-d8f1-4711-8991-bb8f5c92ad7c"
	haloTag.Spec.DisplayName = "Halo"
	haloTag.Spec.Slug = "halo"
	h.store.Create(ctx, haloTag)

	primaryMenu := &model.Menu{}
	primaryMenu.Metadata.Name = "primary"
	primaryMenu.Spec.DisplayName = "主菜单"
	h.store.Create(ctx, primaryMenu)

	homeMenuItem := &model.MenuItem{}
	homeMenuItem.Metadata.Name = "88c3f10b-321c-4092-86a8-70db00251b74"
	homeMenuItem.Spec.DisplayName = "首页"
	homeMenuItem.Spec.Href = "/"
	homeMenuItem.Spec.Priority = 0
	h.store.Create(ctx, homeMenuItem)

	adminRole := &model.Role{}
	adminRole.Metadata.Name = "role-admin"
	adminRole.Spec.DisplayName = "管理员"
	adminRole.Spec.Type = "SYSTEM"
	h.store.Create(ctx, adminRole)
}
