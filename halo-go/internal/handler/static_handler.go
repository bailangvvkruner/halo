package handler

import (
	"bytes"
	"embed"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed all:web
var webFS embed.FS

type StaticHandler struct {
	consoleHTML []byte
	ucHTML      []byte
	workDir     string
}

func NewStaticHandler(workDir string) *StaticHandler {
	consoleData, _ := webFS.ReadFile("web/console.html")
	ucData, _ := webFS.ReadFile("web/uc.html")
	return &StaticHandler{
		consoleHTML: consoleData,
		ucHTML:      ucData,
		workDir:     workDir,
	}
}

func (h *StaticHandler) ServeConsole(c *gin.Context) {
	c.DataFromReader(http.StatusOK, int64(len(h.consoleHTML)), "text/html; charset=utf-8",
		bytes.NewReader(h.consoleHTML), nil)
	c.Header("Cache-Control", "no-store")
}

func (h *StaticHandler) ServeUC(c *gin.Context) {
	c.DataFromReader(http.StatusOK, int64(len(h.ucHTML)), "text/html; charset=utf-8",
		bytes.NewReader(h.ucHTML), nil)
	c.Header("Cache-Control", "no-store")
}

func (h *StaticHandler) UIAssetsFS() http.FileSystem {
	subFS, _ := fs.Sub(webFS, "web/ui-assets")
	return http.FS(subFS)
}

func (h *StaticHandler) ServeUIAssetsDirect(c *gin.Context) {
	filepath := c.Param("filepath")
	if filepath == "" || filepath == "/" {
		c.Status(http.StatusNotFound)
		return
	}
	filepath = path.Clean(filepath)
	if strings.HasPrefix(filepath, "/") {
		filepath = filepath[1:]
	}
	data, err := webFS.ReadFile("web/ui-assets/" + filepath)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	contentType := h.guessContentType(filepath)
	c.DataFromReader(http.StatusOK, int64(len(data)), contentType,
		bytes.NewReader(data), nil)
}

func (h *StaticHandler) guessContentType(name string) string {
	ext := path.Ext(name)
	switch ext {
	case ".ico":
		return "image/x-icon"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".svg":
		return "image/svg+xml"
	case ".css":
		return "text/css; charset=utf-8"
	case ".js", ".mjs":
		return "application/javascript"
	case ".woff2":
		return "font/woff2"
	case ".woff":
		return "font/woff"
	case ".ttf":
		return "font/ttf"
	case ".eot":
		return "application/vnd.ms-fontobject"
	case ".json":
		return "application/json"
	case ".html", ".htm":
		return "text/html; charset=utf-8"
	default:
		return "application/octet-stream"
	}
}

func (h *StaticHandler) ServeThemeAssets(c *gin.Context) {
	themeName := c.Param("themeName")
	filepath := c.Param("filepath")
	if filepath == "" || filepath == "/" {
		c.Status(http.StatusNotFound)
		return
	}
	filepath = path.Clean(filepath)
	if strings.HasPrefix(filepath, "/") {
		filepath = filepath[1:]
	}
	assetPath := path.Join(h.workDir, "themes", themeName, "assets", filepath)
	data, err := os.ReadFile(assetPath)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	contentType := h.guessContentType(filepath)
	c.Header("Cache-Control", "public, max-age=31536000")
	c.DataFromReader(http.StatusOK, int64(len(data)), contentType,
		bytes.NewReader(data), nil)
}

func (h *StaticHandler) IsHTMLRequest(c *gin.Context) bool {
	accept := c.GetHeader("Accept")
	if accept == "" {
		return false
	}
	acceptLower := strings.ToLower(accept)
	types := strings.Split(acceptLower, ",")
	for _, t := range types {
		t = strings.TrimSpace(t)
		if idx := strings.Index(t, ";"); idx >= 0 {
			t = t[:idx]
		}
		t = strings.TrimSpace(t)
		parsed, err := url.Parse(t)
		if err == nil && parsed != nil {
			mimeType := parsed.Path
			if mimeType == "text/html" || mimeType == "*/*" || mimeType == "text/*" {
				return true
			}
		} else if t == "text/html" || t == "*/*" || t == "text/*" {
			return true
		}
	}
	return false
}
