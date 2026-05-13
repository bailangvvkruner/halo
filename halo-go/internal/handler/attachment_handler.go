package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/halo-dev/halo-go/internal/config"
	"github.com/halo-dev/halo-go/internal/model"
	"github.com/halo-dev/halo-go/internal/service"
)

type AttachmentHandler struct {
	svc  service.AttachmentService
	cfg  *config.Config
}

func NewAttachmentHandler(svc service.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{svc: svc}
}

func (h *AttachmentHandler) Create(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请选择要上传的文件", "data": nil})
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	newFileName := fmt.Sprintf("%s%s", uuid.New().String()[:8], ext)
	uploadDir := filepath.Join(h.cfg.WorkDir, "attachments", time.Now().Format("2006/01"))
	os.MkdirAll(uploadDir, 0755)
	filePath := filepath.Join(uploadDir, newFileName)

	if out, err := os.Create(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建文件失败: " + err.Error(), "data": nil})
		return
	} else {
		defer out.Close()
		if _, err := out.ReadFrom(file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存文件失败: " + err.Error(), "data": nil})
			return
		}
	}

	attachment := &model.Attachment{}
	attachment.Metadata.Name = uuid.New().String()
	attachment.Spec.DisplayName = header.Filename
	attachment.Spec.MediaType = header.Header.Get("Content-Type")
	attachment.Spec.Size = header.Size
	attachment.Spec.Path = filePath

	created, err := h.svc.Create(c.Request.Context(), attachment)
	if err != nil {
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建附件记录失败: " + err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "上传成功",
		"data": gin.H{
			"name":     created.Metadata.Name,
			"displayName": created.Spec.DisplayName,
			"path":     created.Spec.Path,
			"url":      fmt.Sprintf("/api/v1alpha1/attachments/%s/download", created.Metadata.Name),
			"size":     created.Spec.Size,
			"mediaType": created.Spec.MediaType,
		},
	})
}

func (h *AttachmentHandler) Get(c *gin.Context) {
	name := c.Param("name")
	attachment, err := h.svc.Get(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "附件不存在: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": attachment})
}

func (h *AttachmentHandler) Update(c *gin.Context) {
	name := c.Param("name")
	var attachment model.Attachment
	if err := c.ShouldBindJSON(&attachment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误: " + err.Error(), "data": nil})
		return
	}
	attachment.Metadata.Name = name
	updated, err := h.svc.Update(c.Request.Context(), &attachment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新附件失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": updated})
}

func (h *AttachmentHandler) Delete(c *gin.Context) {
	name := c.Param("name")
	if err := h.svc.Delete(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除附件失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": nil})
}

func (h *AttachmentHandler) List(c *gin.Context) {
	opts := parseListOptions(c)
	result, err := h.svc.List(c.Request.Context(), &opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询附件列表失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": result})
}
