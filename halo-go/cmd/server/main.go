package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/halo-dev/halo-go/internal/config"
	"github.com/halo-dev/halo-go/internal/data"
	"github.com/halo-dev/halo-go/internal/router"
	"github.com/halo-dev/halo-go/internal/server"
	"github.com/halo-dev/halo-go/internal/service"
	"github.com/halo-dev/halo-go/internal/seed"
)

func main() {
	log.Println("============================================================")
	log.Println("  Halo Go CMS Blog System v0.1.0")
	log.Println("  基于 Go + Gin + GORM + SQLite 构建")
	log.Println("============================================================")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	db, err := data.InitDB(cfg.Database.DataSource)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	store := data.NewExtensionStore(db)
	client := data.NewClient(store)

	postService := service.NewPostService(client)
	categoryService := service.NewCategoryService(client)
	tagService := service.NewTagService(client)
	singlePageService := service.NewSinglePageService(client)
	replyService := service.NewReplyService(client)
	commentService := service.NewCommentService(client, replyService)
	userService := service.NewUserService(client)
	roleService := service.NewRoleService(client)
	pluginService := service.NewPluginService(client)
	themeService := service.NewThemeService(client)
	menuService := service.NewMenuService(client)
	menuItemService := service.NewMenuItemService(client)
	attachmentService := service.NewAttachmentService(client)
	notificationService := service.NewNotificationService(client)
	authService := service.NewAuthService(userService, cfg.JWT.Secret, cfg.JWT.ExpireHours)
	snapshotService := service.NewSnapshotService(client)
	statsService := service.NewStatsService()

	srv := &router.Server{
		PostService:         postService,
		CategoryService:     categoryService,
		TagService:          tagService,
		SinglePageService:   singlePageService,
		CommentService:      commentService,
		ReplyService:        replyService,
		UserService:         userService,
		RoleService:         roleService,
		PluginService:       pluginService,
		ThemeService:        themeService,
		MenuService:         menuService,
		MenuItemService:     menuItemService,
		AttachmentService:   attachmentService,
		NotificationService: notificationService,
		AuthService:         authService,
		SnapshotService:     snapshotService,
		StatsService:        statsService,
	}

	engine, err := server.New(cfg, db)
	if err != nil {
		log.Fatalf("创建服务器引擎失败: %v", err)
	}

	router.RegisterRoutes(engine, srv, cfg, store)

	defaultUser := "admin"
	if err := seed.SeedData(store, defaultUser); err != nil {
		log.Printf("警告: 初始数据播种失败 (可能已存在): %v", err)
	} else {
		log.Println("初始数据播种成功")
	}

	httpServer := &http.Server{
		Addr:    cfg.Addr(),
		Handler: engine,
	}

	go func() {
		log.Printf("Halo Go 服务器启动，监听地址: %s", cfg.Addr())
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器运行错误: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("收到信号: %v，正在优雅关闭服务器...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("服务器强制关闭: %v", err)
	}

	log.Println("服务器已优雅关闭")
}
