package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/halo-dev/halo-go/internal/config"
	"github.com/halo-dev/halo-go/internal/router"
	"gorm.io/gorm"
)

func New(cfg *config.Config, db *gorm.DB) (*gin.Engine, error) {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "接口不存在",
			"data":    nil,
		})
	})
	return r, nil
}

func Run(engine *gin.Engine, cfg *config.Config) error {
	srv := &http.Server{
		Addr:         cfg.Addr(),
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Printf("Halo Go 服务器启动，监听地址: %s", cfg.Addr())
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("服务器启动失败: %w", err)
	}
	return nil
}

func RunWithGracefulShutdown(engine *gin.Engine, cfg *config.Config) error {
	srv := &http.Server{
		Addr:         cfg.Addr(),
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Printf("Halo Go 服务器启动，监听地址: %s", cfg.Addr())
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("服务器错误: %v", err)
		}
	}()
	quit := make(chan struct{})
	go func() {
		time.Sleep(100 * time.Millisecond)
		close(quit)
	}()
	<-quit
	return srv.Shutdown(nil)
}

var _ = router.Server{}
