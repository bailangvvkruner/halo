package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/halo-dev/halo-go/internal/config"
	"github.com/halo-dev/halo-go/internal/extension"
	"github.com/halo-dev/halo-go/internal/model"
	"github.com/halo-dev/halo-go/internal/service"
)

type PluginMarketHandler struct {
	svc  service.PluginService
	cfg  *config.Config
}

func NewPluginMarketHandler(svc service.PluginService) *PluginMarketHandler {
	return &PluginMarketHandler{svc: svc}
}

type PluginManifest struct {
	APIVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind       string `yaml:"kind" json:"kind"`
	Metadata   struct {
		Name string `yaml:"name" json:"name"`
	} `yaml:"metadata" json:"metadata"`
	Spec struct {
		DisplayName    string   `yaml:"displayName" json:"displayName"`
		Description    string   `yaml:"description" json:"description"`
		Version        string   `yaml:"version" json:"version"`
		Author         struct {
			Name  string `yaml:"name" json:"name"`
			Website string `yaml:"website" json:"website"`
		} `yaml:"author" json:"author"`
		Logo           string   `yaml:"logo" json:"logo"`
		Homepage       string   `yaml:"homepage" json:"homepage"`
		Repo           string   `yaml:"repo" json:"repo"`
		Requires       string   `yaml:"requires" json:"requires"`
		Permission      []string `yaml:"permission" json:"permission"`
	} `yaml:"spec" json:"spec"`
}

func (h *PluginMarketHandler) InstallFromURL(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "data": nil})
		return
	}

	tmpDir := filepath.Join(os.TempDir(), "plugin-download")
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	pluginPath := filepath.Join(tmpDir, "plugin.jar")
	out, _ := os.Create(pluginPath)

	resp, err := http.Get(req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "下载失败: " + err.Error(), "data": nil})
		return
	}
	defer resp.Body.Close()

	io.Copy(out, resp.Body)
	out.Close()

	plugin := &model.Plugin{}
	plugin.Metadata.Name = generatePluginName()
	plugin.Spec.DisplayName = filepath.Base(req.URL)
	plugin.Spec.Enabled = false

	created, err := h.svc.Create(c.Request.Context(), plugin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "安装失败: " + err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "插件安装成功",
		"data":    created,
	})
}

func (h *PluginMarketHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请选择插件文件", "data": nil})
		return
	}
	defer file.Close()

	pluginsDir := filepath.Join(h.cfg.WorkDir, "plugins")
	os.MkdirAll(pluginsDir, 0755)

	destPath := filepath.Join(pluginsDir, header.Filename)
	out, _ := os.Create(destPath)
	io.Copy(out, file)
	out.Close()

	plugin := &model.Plugin{}
	plugin.Metadata.Name = generatePluginName()
	plugin.Spec.DisplayName = strings.TrimSuffix(header.Filename, ".jar")
	plugin.Spec.Enabled = false

	created, err := h.svc.Create(c.Request.Context(), plugin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "注册插件失败: " + err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "插件上传成功",
		"data":    created,
	})
}

func (h *PluginMarketHandler) GetPresetPlugins(c *gin.Context) {
	presetPlugins := []PluginManifest{
		{
			APIVersion: "v1alpha1",
			Kind:       "Plugin",
			Spec: struct {
				DisplayName    string   `yaml:"displayName" json:"displayName"`
				Description    string   `yaml:"description" json:"description"`
				Version        string   `yaml:"version" json:"version"`
				Author         struct {
					Name  string `yaml:"name" json:"name"`
					Website string `yaml:"website" json:"website"`
				} `yaml:"author" json:"author"`
				Logo           string   `yaml:"logo" json:"logo"`
				Homepage       string   `yaml:"homepage" json:"homepage"`
				Repo           string   `yaml:"repo" json:"repo"`
				Requires       string   `yaml:"requires" json:"requires"`
				Permission      []string `yaml:"permission" json:"permission"`
			}{
				DisplayName: "Sitemap",
				Description:  "自动生成站点地图",
				Version:      "1.0.0",
				Author: struct {
					Name  string `yaml:"name" json:"name"`
					Website string `yaml:"website" json:"website"`
				}{
					Name:  "Halo Team",
					Website: "https://www.halo.run",
				},
				Homepage: "https://www.halo.run/store/apps/app-SfwaE",
			},
		},
		{
			APIVersion: "v1alpha1",
			Kind:       "Plugin",
			Spec: struct {
				DisplayName    string   `yaml:"displayName" json:"displayName"`
				Description    string   `yaml:"description" json:"description"`
				Version        string   `yaml:"version" json:"version"`
				Author         struct {
					Name  string `yaml:"name" json:"name"`
					Website string `yaml:"website" json:"website"`
				} `yaml:"author" json:"author"`
				Logo           string   `yaml:"logo" json:"logo"`
				Homepage       string   `yaml:"homepage" json:"homepage"`
				Repo           string   `yaml:"repo" json:"repo"`
				Requires       string   `yaml:"requires" json:"requires"`
				Permission      []string `yaml:"permission" json:"permission"`
			}{
				DisplayName: "Search Plugin",
				Description:  "全文搜索增强",
				Version:      "1.2.0",
				Author: struct {
					Name  string `yaml:"name" json:"name"`
					Website string `yaml:"website" json:"website"`
				}{
					Name:  "Halo Team",
					Website: "https://www.halo.run",
				},
				Homepage: "https://www.halo.run/store/apps/app-SEARCH",
			},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    presetPlugins,
	})
}

func (h *PluginMarketHandler) CheckUpdate(c *gin.Context) {
	name := c.Param("name")

	result, _ := h.svc.List(c.Request.Context(), &extension.ListOptions{Size: 0})
	for _, item := range result.Items {
		if plugin, ok := item.(*model.Plugin); ok && plugin.Metadata.Name == name {
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "ok",
				"data": gin.H{
					"currentVersion": plugin.Spec.Version,
					"latestVersion":  plugin.Spec.Version,
					"hasUpdate":     false,
				},
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "插件不存在", "data": nil})
}

func generatePluginName() string {
	return fmt.Sprintf("plugin-%d", len(filepath.Join(os.TempDir())))
}
