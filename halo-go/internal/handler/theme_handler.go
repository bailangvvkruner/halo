package handler

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"

	"github.com/halo-dev/halo-go/internal/config"
	"github.com/halo-dev/halo-go/internal/model"
	"github.com/halo-dev/halo-go/internal/service"
)

type ThemeHandler struct {
	svc  service.ThemeService
	cfg  *config.Config
}

func NewThemeHandler(svc service.ThemeService) *ThemeHandler {
	return &ThemeHandler{svc: svc}
}

func (h *ThemeHandler) Create(c *gin.Context) {
	var theme model.Theme
	if err := c.ShouldBindJSON(&theme); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误: " + err.Error(), "data": nil})
		return
	}
	created, err := h.svc.Create(c.Request.Context(), &theme)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建主题失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": created})
}

func (h *ThemeHandler) Install(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请选择主题文件", "data": nil})
		return
	}
	defer file.Close()

	tmpDir := filepath.Join(os.TempDir(), fmt.Sprintf("theme-upload-%d", header.Size))
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	zipPath := filepath.Join(tmpDir, "theme.zip")
	out, _ := os.Create(zipPath)
	io.Copy(out, file)
	out.Close()

	themesDir := filepath.Join(h.cfg.WorkDir, "themes")
	os.MkdirAll(themesDir, 0755)

	if err := unzip(zipPath, themesDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "解压主题失败: " + err.Error(), "data": nil})
		return
	}

	themeName := strings.TrimSuffix(header.Filename, ".zip")
	themeSettingsPath := filepath.Join(themesDir, themeName, "theme.yaml")
	if _, err := os.Stat(themeSettingsPath); os.IsNotExist(err) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的主题文件：缺少 theme.yaml", "data": nil})
		return
	}

	data, _ := os.ReadFile(themeSettingsPath)
	var themeSpec model.ThemeSpec
	yaml.Unmarshal(data, &themeSpec)

	theme := &model.Theme{}
	theme.Metadata.Name = generateThemeName()
	theme.Spec = themeSpec
	theme.Spec.Installed = true

	created, err := h.svc.Create(c.Request.Context(), theme)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "注册主题失败: " + err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "主题安装成功",
		"data":    created,
	})
}

func (h *ThemeHandler) Get(c *gin.Context) {
	name := c.Param("name")
	theme, err := h.svc.Get(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "主题不存在: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": theme})
}

func (h *ThemeHandler) Update(c *gin.Context) {
	name := c.Param("name")
	var theme model.Theme
	if err := c.ShouldBindJSON(&theme); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误: " + err.Error(), "data": nil})
		return
	}
	theme.Metadata.Name = name
	updated, err := h.svc.Update(c.Request.Context(), &theme)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新主题失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": updated})
}

func (h *ThemeHandler) Delete(c *gin.Context) {
	name := c.Param("name")
	if err := h.svc.Delete(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除主题失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": nil})
}

func (h *ThemeHandler) List(c *gin.Context) {
	opts := parseListOptions(c)
	result, err := h.svc.List(c.Request.Context(), &opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询主题列表失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": result})
}

func (h *ThemeHandler) Activate(c *gin.Context) {
	name := c.Param("name")
	if err := h.svc.Activate(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "激活主题失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": nil})
}

func (h *ThemeHandler) GetConfig(c *gin.Context) {
	name := c.Param("name")
	configPath := filepath.Join(h.cfg.WorkDir, "themes", name, "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": map[string]interface{}{}})
		return
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "读取配置失败", "data": nil})
		return
	}
	var configData map[string]interface{}
	json.Unmarshal(data, &configData)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": configData})
}

func (h *ThemeHandler) SaveConfig(c *gin.Context) {
	name := c.Param("name")
	var configData map[string]interface{}
	if err := c.ShouldBindJSON(&configData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "data": nil})
		return
	}
	configPath := filepath.Join(h.cfg.WorkDir, "themes", name, "config.yaml")
	data, _ := json.MarshalIndent(configData, "", "  ")
	os.WriteFile(configPath, data, 0644)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "配置保存成功", "data": nil})
}

func (h *ThemeHandler) GetTemplates(c *gin.Context) {
	name := c.Param("name")
	templateDir := filepath.Join(h.cfg.WorkDir, "themes", name, "templates")
	files, err := os.ReadDir(templateDir)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": []string{}})
		return
	}
	templates := make([]string, 0)
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".html") {
			templates = append(templates, f.Name())
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": templates})
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		path := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			continue
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(path), 0755)
		out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		_, err = io.Copy(out, rc)
		out.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func generateThemeName() string {
	return fmt.Sprintf("theme-%d", len(filepath.Join(os.TempDir())))
}
