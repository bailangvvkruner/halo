package router

import (
	"github.com/gin-gonic/gin"

	"github.com/halo-dev/halo-go/internal/config"
	"github.com/halo-dev/halo-go/internal/data"
	"github.com/halo-dev/halo-go/internal/extension"
	"github.com/halo-dev/halo-go/internal/handler"
	"github.com/halo-dev/halo-go/internal/middleware"
	"github.com/halo-dev/halo-go/internal/service"
)

type Server struct {
	PostService        service.PostService
	CategoryService    service.CategoryService
	TagService         service.TagService
	SinglePageService  service.SinglePageService
	CommentService     service.CommentService
	ReplyService       service.ReplyService
	UserService        service.UserService
	RoleService        service.RoleService
	PluginService      service.PluginService
	ThemeService       service.ThemeService
	MenuService        service.MenuService
	MenuItemService    service.MenuItemService
	AttachmentService  service.AttachmentService
	NotificationService service.NotificationService
	AuthService        service.AuthService
	SnapshotService    service.SnapshotService
	StatsService       service.StatsService
}

func RegisterRoutes(r *gin.Engine, srv *Server, cfg *config.Config, store *data.ExtensionStore) {
	authMiddleware := middleware.AuthMiddleware(cfg.JWT.Secret)
	apiV1 := r.Group("/api/v1alpha1")
	apiV1.Use(authMiddleware)
	{
		postHandler := handler.NewPostHandler(srv.PostService)
		posts := apiV1.Group("/posts")
		{
			posts.POST("", postHandler.Create)
			posts.GET("", postHandler.List)
			posts.GET("/:name", postHandler.Get)
			posts.PUT("/:name", postHandler.Update)
			posts.DELETE("/:name", postHandler.Delete)
			posts.PUT("/:name/publish", postHandler.Publish)
			posts.PUT("/:name/unpublish", postHandler.Unpublish)
			posts.PUT("/:name/trash", postHandler.Trash)
			posts.PUT("/:name/restore", postHandler.Restore)
		}
		categoryHandler := handler.NewCategoryHandler(srv.CategoryService)
		categories := apiV1.Group("/categories")
		{
			categories.POST("", categoryHandler.Create)
			categories.GET("", categoryHandler.List)
			categories.GET("/:name", categoryHandler.Get)
			categories.PUT("/:name", categoryHandler.Update)
			categories.DELETE("/:name", categoryHandler.Delete)
		}
		tagHandler := handler.NewTagHandler(srv.TagService)
		tags := apiV1.Group("/tags")
		{
			tags.POST("", tagHandler.Create)
			tags.GET("", tagHandler.List)
			tags.GET("/:name", tagHandler.Get)
			tags.PUT("/:name", tagHandler.Update)
			tags.DELETE("/:name", tagHandler.Delete)
		}
		commentHandler := handler.NewCommentHandler(srv.CommentService)
		comments := apiV1.Group("/comments")
		{
			comments.POST("", commentHandler.Create)
			comments.GET("", commentHandler.List)
			comments.GET("/:name", commentHandler.Get)
			comments.PUT("/:name", commentHandler.Update)
			comments.DELETE("/:name", commentHandler.Delete)
			comments.PUT("/:name/approve", commentHandler.Approve)
			comments.POST("/:name/reply", commentHandler.Reply)
		}
		replyHandler := handler.NewReplyHandler(srv.ReplyService)
		replies := apiV1.Group("/replies")
		{
			replies.POST("", replyHandler.Create)
			replies.GET("", replyHandler.List)
			replies.GET("/:name", replyHandler.Get)
			replies.PUT("/:name", replyHandler.Update)
			replies.DELETE("/:name", replyHandler.Delete)
		}
		userHandler := handler.NewUserHandler(srv.UserService)
		users := apiV1.Group("/users")
		{
			users.POST("", userHandler.Create)
			users.GET("", userHandler.List)
			users.GET("/:name", userHandler.Get)
			users.PUT("/:name", userHandler.Update)
			users.DELETE("/:name", userHandler.Delete)
			users.PUT("/:name/password", userHandler.ChangePassword)
			users.PUT("/:name/profile", userHandler.UpdateProfile)
		}
		roleHandler := handler.NewRoleHandler(srv.RoleService)
		roles := apiV1.Group("/roles")
		{
			roles.POST("", roleHandler.Create)
			roles.GET("", roleHandler.List)
			roles.GET("/:name", roleHandler.Get)
			roles.PUT("/:name", roleHandler.Update)
			roles.DELETE("/:name", roleHandler.Delete)
		}
		pluginHandler := handler.NewPluginHandler(srv.PluginService)
		plugins := apiV1.Group("/plugins")
		{
			plugins.POST("", pluginHandler.Create)
			plugins.GET("", pluginHandler.List)
			plugins.GET("/:name", pluginHandler.Get)
			plugins.PUT("/:name", pluginHandler.Update)
			plugins.DELETE("/:name", pluginHandler.Delete)
			plugins.PUT("/:name/enable", pluginHandler.Enable)
			plugins.PUT("/:name/disable", pluginHandler.Disable)
		}
		themeHandler := handler.NewThemeHandler(srv.ThemeService)
		themes := apiV1.Group("/themes")
		{
			themes.POST("", themeHandler.Create)
			themes.GET("", themeHandler.List)
			themes.GET("/:name", themeHandler.Get)
			themes.PUT("/:name", themeHandler.Update)
			themes.DELETE("/:name", themeHandler.Delete)
			themes.PUT("/:name/activate", themeHandler.Activate)
		}
		menuHandler := handler.NewMenuHandler(srv.MenuService)
		menus := apiV1.Group("/menus")
		{
			menus.POST("", menuHandler.Create)
			menus.GET("", menuHandler.List)
			menus.GET("/:name", menuHandler.Get)
			menus.PUT("/:name", menuHandler.Update)
			menus.DELETE("/:name", menuHandler.Delete)
		}
		menuItemHandler := handler.NewMenuItemHandler(srv.MenuItemService)
		menuItems := apiV1.Group("/menuitems")
		{
			menuItems.POST("", menuItemHandler.Create)
			menuItems.GET("", menuItemHandler.List)
			menuItems.GET("/:name", menuItemHandler.Get)
			menuItems.PUT("/:name", menuItemHandler.Update)
			menuItems.DELETE("/:name", menuItemHandler.Delete)
		}
		attachmentHandler := handler.NewAttachmentHandler(srv.AttachmentService)
		attachments := apiV1.Group("/attachments")
		{
			attachments.POST("", attachmentHandler.Create)
			attachments.GET("", attachmentHandler.List)
			attachments.GET("/:name", attachmentHandler.Get)
			attachments.PUT("/:name", attachmentHandler.Update)
			attachments.DELETE("/:name", attachmentHandler.Delete)
		}
		notificationHandler := handler.NewNotificationHandler(srv.NotificationService)
		notifications := apiV1.Group("/notifications")
		{
			notifications.POST("", notificationHandler.Create)
			notifications.GET("", notificationHandler.List)
			notifications.GET("/:name", notificationHandler.Get)
			notifications.PUT("/:name", notificationHandler.Update)
			notifications.DELETE("/:name", notificationHandler.Delete)
			notifications.PUT("/:name/read", notificationHandler.MarkRead)
		}
		singlePageHandler := handler.NewSinglePageHandler(srv.SinglePageService)
		singlePages := apiV1.Group("/singlepages")
		{
			singlePages.POST("", singlePageHandler.Create)
			singlePages.GET("", singlePageHandler.List)
			singlePages.GET("/:name", singlePageHandler.Get)
			singlePages.PUT("/:name", singlePageHandler.Update)
			singlePages.DELETE("/:name", singlePageHandler.Delete)
			singlePages.PUT("/:name/publish", singlePageHandler.Publish)
			singlePages.PUT("/:name/unpublish", singlePageHandler.Unpublish)
			singlePages.PUT("/:name/trash", singlePageHandler.Trash)
			singlePages.PUT("/:name/restore", singlePageHandler.Restore)
		}
		statsHandler := handler.NewStatsHandler(srv.StatsService)
		stats := apiV1.Group("/stats")
		{
			stats.GET("", statsHandler.GetAllStats)
			stats.GET("/:name", statsHandler.GetVisitCount)
			stats.POST("/:name/visit", statsHandler.IncrVisit)
		}
	}
	publicAPI := r.Group("/api/public")
	{
		publicPosts := publicAPI.Group("/posts")
		{
			publicPosts.GET("", func(c *gin.Context) {
				postHandler := handler.NewPostHandler(srv.PostService)
				postHandler.List(c)
			})
			publicPosts.GET("/:name", func(c *gin.Context) {
				postHandler := handler.NewPostHandler(srv.PostService)
				postHandler.Get(c)
			})
		}
		publicCategories := publicAPI.Group("/categories")
		{
			publicCategories.GET("", func(c *gin.Context) {
				categoryHandler := handler.NewCategoryHandler(srv.CategoryService)
				categoryHandler.List(c)
			})
		}
		publicTags := publicAPI.Group("/tags")
		{
			publicTags.GET("", func(c *gin.Context) {
				tagHandler := handler.NewTagHandler(srv.TagService)
				tagHandler.List(c)
			})
		}
	}
	actuatorHandler := handler.NewActuatorHandler()
	actuator := r.Group("/actuator")
	{
		actuator.GET("/health", actuatorHandler.Health)
		actuator.GET("/info", actuatorHandler.Info)
	}
	authHandler := handler.NewAuthHandler(srv.AuthService)
	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/refresh-token", authHandler.RefreshToken)
		auth.POST("/password-reset", authHandler.PasswordReset)
		auth.GET("/me", authMiddleware, authHandler.CurrentUser)
	}

	setupHandler := handler.NewSetupHandler(store)
	setup := r.Group("/api/v1alpha1")
	{
		setup.GET("/setup", setupHandler.GetStatus)
		setup.POST("/setup", setupHandler.DoSetup)
	}

	staticHandler := handler.NewStaticHandler()
	r.GET("/console/*filepath", staticHandler.ServeConsole)
	r.GET("/uc/*filepath", staticHandler.ServeUC)
	r.StaticFS("/ui-assets", staticHandler.UIAssetsFS())
}

var _ = extension.GVK{}
