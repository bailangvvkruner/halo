package config

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	JWT      JWTConfig      `json:"jwt"`
	WorkDir  string         `json:"workDir"`
}

type ServerConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

type DatabaseConfig struct {
	DataSource string `json:"dataSource"`
}

type JWTConfig struct {
	Secret      string `json:"secret"`
	ExpireHours int    `json:"expireHours"`
}

var (
	serverPort  = flag.Int("server.port", 0, "服务端口 (同 Java Halo: --server.port)")
	workDir     = flag.String("halo.work-dir", "", "工作目录 (同 Java Halo: --halo.work-dir)")
)

func Load() (*Config, error) {
	if !flag.Parsed() {
		flag.Parse()
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: 8090,
			Host: "0.0.0.0",
		},
		Database: DatabaseConfig{
			DataSource: "",
		},
		JWT: JWTConfig{
			Secret:      "halo-default-secret-change-me",
			ExpireHours: 24,
		},
		WorkDir: "",
	}

	if *serverPort > 0 {
		cfg.Server.Port = *serverPort
	}
	if v := os.Getenv("HALO_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil && port > 0 {
			cfg.Server.Port = port
		}
	}
	if v := os.Getenv("SERVER_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil && port > 0 {
			cfg.Server.Port = port
		}
	}

	if *workDir != "" {
		cfg.WorkDir = *workDir
	}
	if v := os.Getenv("HALO_WORK_DIR"); v != "" {
		cfg.WorkDir = v
	}
	if cfg.WorkDir == "" {
		cfg.WorkDir = "./data"
	}

	if cfg.Database.DataSource == "" {
		cfg.Database.DataSource = cfg.WorkDir + "/halo.db"
	}
	if v := os.Getenv("HALO_DB_DATASOURCE"); v != "" {
		cfg.Database.DataSource = v
	}

	if v := os.Getenv("HALO_JWT_SECRET"); v != "" {
		cfg.JWT.Secret = v
	}

	if err := os.MkdirAll(cfg.WorkDir, 0755); err != nil {
		return nil, fmt.Errorf("创建工作目录失败 %s: %w", cfg.WorkDir, err)
	}

	return cfg, nil
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func (c *Config) Shutdown(engine *gin.Engine, ctx context.Context) error {
	srv := &http.Server{
		Addr:    c.Addr(),
		Handler: engine,
	}
	log.Println("正在关闭 HTTP 服务器...")
	return srv.Shutdown(ctx)
}