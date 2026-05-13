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
	"time"

	"github.com/gin-gonic/gin"

	"github.com/halo-dev/halo-go/internal/config"
	"github.com/halo-dev/halo-go/internal/data"
	"github.com/halo-dev/halo-go/internal/extension"
)

type BackupHandler struct {
	store *data.ExtensionStore
	cfg   *config.Config
}

func NewBackupHandler(store *data.ExtensionStore, cfg *config.Config) *BackupHandler {
	return &BackupHandler{store: store, cfg: cfg}
}

func (h *BackupHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	result, err := h.store.List(ctx, &extension.ListOptions{Size: 0})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取数据失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	backupData := BackupManifest{
		Version:     "1.0",
		CreatedAt:   time.Now().Format(time.RFC3339),
		HaloVersion: "2.0.0-go",
		Data:        result.Items,
	}

	jsonData, _ := json.MarshalIndent(backupData, "", "  ")

	backupDir := filepath.Join(h.cfg.WorkDir, "backups")
	os.MkdirAll(backupDir, 0755)
	filename := fmt.Sprintf("halo-backup-%s.zip", time.Now().Format("20060102-150405"))
	backupPath := filepath.Join(backupDir, filename)

	zipFile, err := os.Create(backupPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建备份文件失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	writer, _ := zipWriter.Create("manifest.json")
	writer.Write(jsonData)

	dbPath := filepath.Join(h.cfg.WorkDir, "data", "halo-go.db")
	if _, err := os.Stat(dbPath); err == nil {
		h.addToZip(zipWriter, dbPath, "halo-go.db")
	}

	attachmentsDir := filepath.Join(h.cfg.WorkDir, "attachments")
	if _, err := os.Stat(attachmentsDir); err == nil {
		h.addDirToZip(zipWriter, attachmentsDir, "attachments")
	}

	zipWriter.Close()

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "备份创建成功",
		"data": gin.H{
			"filename": filename,
			"path":     backupPath,
			"size":     fileInfo(backupPath),
		},
	})
}

func (h *BackupHandler) List(c *gin.Context) {
	backupDir := filepath.Join(h.cfg.WorkDir, "backups")
	files, err := os.ReadDir(backupDir)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": []string{}})
		return
	}

	backups := make([]BackupInfo, 0)
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".zip") {
			info, _ := f.Info()
			backups = append(backups, BackupInfo{
				Filename: f.Name(),
				Size:     info.Size(),
				Created:  info.ModTime().Format(time.RFC3339),
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    backups,
	})
}

func (h *BackupHandler) Download(c *gin.Context) {
	filename := c.Param("filename")
	backupPath := filepath.Join(h.cfg.WorkDir, "backups", filename)

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "备份文件不存在", "data": nil})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.File(backupPath)
}

func (h *BackupHandler) Restore(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请选择备份文件", "data": nil})
		return
	}
	defer file.Close()

	tmpPath := filepath.Join(os.TempDir(), header.Filename)
	out, _ := os.Create(tmpPath)
	io.Copy(out, file)
	out.Close()

	reader, err := zip.OpenReader(tmpPath)
	if err != nil {
		os.Remove(tmpPath)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的备份文件", "data": nil})
		return
	}
	defer reader.Close()

	var manifest BackupManifest
	for _, f := range reader.File {
		if f.Name == "manifest.json" {
			rc, _ := f.Open()
			data, _ := io.ReadAll(rc)
			json.Unmarshal(data, &manifest)
			rc.Close()
			break
		}
	}

	if manifest.Version == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的备份格式", "data": nil})
		return
	}

	dbDest := filepath.Join(h.cfg.WorkDir, "data", "halo-go.db")
	for _, f := range reader.File {
		if f.Name == "halo-go.db" {
			rc, _ := f.Open()
			out, _ := os.Create(dbDest + ".restoring")
			io.Copy(out, rc)
			rc.Close()
			out.Close()
			os.Rename(dbDest+".restoring", dbDest)
		} else if strings.HasPrefix(f.Name, "attachments/") {
			destPath := filepath.Join(h.cfg.WorkDir, f.Name)
			os.MkdirAll(filepath.Dir(destPath), 0755)
			rc, _ := f.Open()
			out, _ := os.Create(destPath)
			io.Copy(out, rc)
			rc.Close()
			out.Close()
		}
	}

	os.Remove(tmpPath)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "恢复成功，请重启服务以加载数据",
		"data":    nil,
	})
}

func (h *BackupHandler) Delete(c *gin.Context) {
	filename := c.Param("filename")
	backupPath := filepath.Join(h.cfg.WorkDir, "backups", filename)
	if err := os.Remove(backupPath); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "备份文件不存在", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}

type BackupManifest struct {
	Version     string                `json:"version"`
	CreatedAt   string                `json:"createdAt"`
	HaloVersion string                `json:"haloVersion"`
	Data        []extension.Extension `json:"data"`
}

type BackupInfo struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Created  string `json:"created"`
}

func (h *BackupHandler) addToZip(z *zip.Writer, path, name string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	info, _ := file.Stat()
	header, _ := zip.FileInfoHeader(info)
	header.Name = name
	header.Method = zip.Deflate
	writer, err := z.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	return err
}

func (h *BackupHandler) addDirToZip(z *zip.Writer, baseDir, prefix string) error {
	filepath.Walk(baseDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		relPath, _ := filepath.Rel(baseDir, filePath)
		return h.addToZip(z, filePath, filepath.Join(prefix, relPath))
	})
	return nil
}

func fileInfo(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}
