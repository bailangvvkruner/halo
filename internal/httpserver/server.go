package httpserver

import (
	"halo/internal/middleware"
	"halo/internal/model"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"halo/internal/config"
	"halo/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
}

func New(cfg config.Config, services *service.Container) *Server {
	g := gin.Default()
	g.Use(cors.Default())

	registerRoutes(g, cfg, services)

	return &Server{engine: g}
}

func (s *Server) Run(addr string) error {
	return s.engine.Run(addr)
}

func registerRoutes(g *gin.Engine, cfg config.Config, services *service.Container) {
	authMw := middleware.AuthMiddleware(services.Auth)
	permissionMw := middleware.PermissionMiddleware(services.Roles)
	api := g.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		api.GET("/dashboard/stats", func(c *gin.Context) {
			initialized, err := services.Setup.IsInitialized()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			if !initialized {
				c.JSON(http.StatusOK, gin.H{"initialized": false})
				return
			}
			stats, err := services.Dashboard.Stats(services)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, stats)
		})

		api.GET("/setup/status", func(c *gin.Context) {
			initialized, err := services.Setup.IsInitialized()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"initialized": initialized})
		})

		api.POST("/setup", func(c *gin.Context) {
			var payload service.SetupPayload
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Setup.Setup(payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			c.Status(http.StatusNoContent)
		})

		api.POST("/login", func(c *gin.Context) {
			var payload struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			token, user, err := services.Auth.Login(payload.Username, payload.Password)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
		})

		api.POST("/register", func(c *gin.Context) {
			var payload model.UserRegistration
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Registration.Register(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"token": payload.Token, "message": "registration created"})
		})

		api.GET("/register/verify", func(c *gin.Context) {
			token := c.Query("token")
			if token == "" {
				c.JSON(http.StatusBadRequest, gin.H{"message": "token required"})
				return
			}
			if err := services.Registration.Verify(token); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "verified"})
		})

		api.GET("/posts", func(c *gin.Context) {
			posts, err := services.Posts.List()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, posts)
		})

		api.POST("/posts", authMw, permissionMw, func(c *gin.Context) {
			var payload struct {
				Title     string `json:"title"`
				Slug      string `json:"slug"`
				Content   string `json:"content"`
				Excerpt   string `json:"excerpt"`
				Published bool   `json:"published"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			post := toPost(payload.Title, payload.Slug, payload.Content, payload.Excerpt, payload.Published)
			if err := services.Posts.Create(post); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, post)
		})

		api.GET("/posts/:id", func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			post, err := services.Posts.Get(id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, post)
		})

		api.PUT("/posts/:id", authMw, permissionMw, func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			var payload struct {
				Title     string `json:"title"`
				Slug      string `json:"slug"`
				Content   string `json:"content"`
				Excerpt   string `json:"excerpt"`
				Published bool   `json:"published"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			post := toPost(payload.Title, payload.Slug, payload.Content, payload.Excerpt, payload.Published)
			post.ID = id
			if err := services.Posts.Update(id, post); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, post)
		})

		api.DELETE("/posts/:id", authMw, permissionMw, func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Posts.Delete(id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.Status(http.StatusNoContent)
		})

		api.GET("/comments", func(c *gin.Context) {
			items, err := services.Comments.List()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, items)
		})

		api.GET("/comments/:id", func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			item, err := services.Comments.Get(id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, item)
		})

		api.POST("/comments", func(c *gin.Context) {
			var payload model.Comment
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Comments.Create(&payload); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, payload)
		})

		api.GET("/comments/:id/replies", func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			replies, err := services.Replies.ListByComment(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, replies)
		})

		api.POST("/comments/:id/replies", func(c *gin.Context) {
			commentID, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			var payload model.Reply
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			payload.CommentID = commentID
			if err := services.Replies.Create(&payload); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, payload)
		})

		api.POST("/replies/:id/approve", authMw, permissionMw, func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Replies.Approve(id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "reply approved"})
		})

		api.POST("/replies/:id/reject", authMw, permissionMw, func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Replies.Reject(id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "reply rejected"})
		})

		api.DELETE("/comments/:id", authMw, permissionMw, func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Comments.Delete(id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.Status(http.StatusNoContent)
		})

		api.GET("/pages", func(c *gin.Context) {
			items, err := services.Pages.List()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, items)
		})

		api.GET("/pages/:id", func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			item, err := services.Pages.Get(id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, item)
		})

		api.POST("/pages", authMw, permissionMw, func(c *gin.Context) {
			var payload model.Page
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Pages.Create(&payload); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, payload)
		})

		api.PUT("/pages/:id", authMw, permissionMw, func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			var payload model.Page
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Pages.Update(id, &payload); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, payload)
		})

		api.DELETE("/pages/:id", authMw, permissionMw, func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Pages.Delete(id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.Status(http.StatusNoContent)
		})

		api.GET("/categories", func(c *gin.Context) {
			items, err := services.Categories.List()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, items)
		})

		api.GET("/categories/:id", func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			item, err := services.Categories.Get(id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, item)
		})

		api.POST("/categories", authMw, permissionMw, func(c *gin.Context) {
			var payload model.Category
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Categories.Create(&payload); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, payload)
		})

		api.PUT("/categories/:id", authMw, permissionMw, func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			var payload model.Category
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Categories.Update(id, &payload); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, payload)
		})

		api.DELETE("/categories/:id", authMw, permissionMw, func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Categories.Delete(id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.Status(http.StatusNoContent)
		})

		api.GET("/tags", func(c *gin.Context) {
			items, err := services.Tags.List()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, items)
		})

		api.GET("/tags/:id", func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			item, err := services.Tags.Get(id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, item)
		})

		api.POST("/tags", authMw, permissionMw, func(c *gin.Context) {
			var payload model.Tag
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Tags.Create(&payload); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, payload)
		})

		api.PUT("/tags/:id", authMw, permissionMw, func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			var payload model.Tag
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Tags.Update(id, &payload); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, payload)
		})

		api.DELETE("/tags/:id", authMw, permissionMw, func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Tags.Delete(id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.Status(http.StatusNoContent)
		})

		api.GET("/menus", func(c *gin.Context) {
			items, err := services.Menus.List()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, items)
		})

		api.GET("/menus/:id", func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			item, err := services.Menus.Get(id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, item)
		})

		api.POST("/menus", authMw, permissionMw, func(c *gin.Context) {
			var payload model.Menu
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Menus.Create(&payload); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, payload)
		})

		api.PUT("/menus/:id", authMw, permissionMw, func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			var payload model.Menu
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Menus.Update(id, &payload); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, payload)
		})

		api.DELETE("/menus/:id", authMw, permissionMw, func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Menus.Delete(id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.Status(http.StatusNoContent)
		})

		api.GET("/themes", func(c *gin.Context) {
			items, err := services.Themes.List()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, items)
		})

		api.GET("/plugins", func(c *gin.Context) {
			items, err := services.Plugins.List()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, items)
		})

		api.GET("/plugins/:id", func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			item, err := services.Plugins.Get(id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, item)
		})

		api.POST("/plugins/:id/enable", authMw, permissionMw, func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Plugins.Toggle(id, true); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "plugin enabled"})
		})

		api.POST("/plugins/:id/disable", authMw, permissionMw, func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Plugins.Toggle(id, false); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "plugin disabled"})
		})

		api.GET("/attachments", func(c *gin.Context) {
			items, err := services.Attachments.List()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, items)
		})

		api.POST("/attachments/upload", authMw, permissionMw, func(c *gin.Context) {
			file, err := c.FormFile("file")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			uploadDir := filepath.Join(cfg.WorkDir, "attachments")
			if err := os.MkdirAll(uploadDir, 0o755); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}

			filename := filepath.Base(file.Filename)
			savePath := filepath.Join(uploadDir, filename)
			if err := c.SaveUploadedFile(file, savePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}

			attachment := &model.Attachment{
				Filename: filename,
				Path:     savePath,
				Size:     file.Size,
			}
			if err := services.Attachments.Create(attachment); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, attachment)
		})

		api.GET("/backups", func(c *gin.Context) {
			items, err := services.Backups.List()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, items)
		})

		api.POST("/backups", authMw, permissionMw, func(c *gin.Context) {
			backup, err := services.Backups.Create(cfg.WorkDir)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, backup)
		})

		api.GET("/backups/:id/download", authMw, permissionMw, func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			backup, err := services.Backups.Get(id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
				return
			}
			c.FileAttachment(filepath.Join(cfg.WorkDir, "backups", backup.Filename), backup.Filename)
		})

		api.GET("/search", func(c *gin.Context) {
			keyword := c.Query("keyword")
			if keyword == "" {
				c.JSON(http.StatusOK, gin.H{"posts": []model.Post{}, "pages": []model.Page{}})
				return
			}

			var posts []model.Post
			if err := services.Posts.Search(keyword, &posts); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}

			var pages []model.Page
			if err := services.Pages.Search(keyword, &pages); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"posts": posts, "pages": pages})
		})

		api.GET("/settings", func(c *gin.Context) {
			items, err := services.Settings.List()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, items)
		})

		api.GET("/settings/:id", func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			item, err := services.Settings.Get(id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, item)
		})

		api.PUT("/settings/:id", authMw, permissionMw, func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			var payload model.Setting
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Settings.Update(id, &payload); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, payload)
		})

		api.GET("/users", func(c *gin.Context) {
			items, err := services.Users.FindAll()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, items)
		})

		api.GET("/users/:id", func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			item, err := services.Users.Get(id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, item)
		})

		api.POST("/users", authMw, permissionMw, func(c *gin.Context) {
			var payload model.User
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Users.Create(&payload); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, payload)
		})

		api.PUT("/users/:id", authMw, permissionMw, func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			var payload model.User
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Users.Update(id, &payload); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, payload)
		})

		api.DELETE("/users/:id", authMw, permissionMw, func(c *gin.Context) {
			id, err := parseUintParam(c, "id")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if err := services.Users.Delete(id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.Status(http.StatusNoContent)
		})
	}

	staticDir := filepath.Join("web", "dist")
	g.GET("/system/setup", func(c *gin.Context) {
		c.File(filepath.Join(staticDir, "index.html"))
	})
	g.Static("/assets", filepath.Join(staticDir, "assets"))
	g.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
			return
		}
		if c.Request.URL.Path == "/" {
			initialized, err := services.Setup.IsInitialized()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			if !initialized {
				c.Redirect(http.StatusFound, "/system/setup")
				return
			}
		}
		c.File(filepath.Join(staticDir, "index.html"))
	})
}
