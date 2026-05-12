package config

import (
	"context"
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

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: 8090,
			Host: "0.0.0.0",
		},
		Database: DatabaseConfig{
			DataSource: "./data/halo.db",
		},
		JWT: JWTConfig{
			Secret:      "halo-default-secret-change-me",
			ExpireHours: 24,
		},
		WorkDir: "./data",
	}
	if v := os.Getenv("HALO_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil && port > 0 {
			cfg.Server.Port = port
		}
	}
	if v := os.Getenv("HALO_DB_DATASOURCE"); v != "" {
		cfg.Database.DataSource = v
	}
	if v := os.Getenv("HALO_JWT_SECRET"); v != "" {
		cfg.JWT.Secret = v
	}
	if v := os.Getenv("HALO_WORK_DIR"); v != "" {
		cfg.WorkDir = v
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
