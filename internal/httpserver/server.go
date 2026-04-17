package httpserver

import (
	"net/http"
	"path/filepath"
	"strings"

	"halo/internal/config"
	"halo/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
}

func New(cfg config.Config, services *service.Container) *Server {
	g := gin.Default()
	g.Use(cors.Default())
	g.Use(sessions.Sessions("halo", cookie.NewStore([]byte(cfg.SessionKey))))

	registerRoutes(g, cfg, services)

	return &Server{engine: g}
}

func (s *Server) Run(addr string) error {
	return s.engine.Run(addr)
}

func registerRoutes(g *gin.Engine, cfg config.Config, services *service.Container) {
	api := g.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
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

		api.GET("/posts", func(c *gin.Context) {
			posts, err := services.Posts.List()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, posts)
		})

		api.POST("/posts", func(c *gin.Context) {
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

		api.GET("/comments", func(c *gin.Context) {
			items, err := services.Comments.List()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, items)
		})

		api.GET("/pages", func(c *gin.Context) {
			items, err := services.Pages.List()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, items)
		})

		api.GET("/categories", func(c *gin.Context) {
			items, err := services.Categories.List()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, items)
		})

		api.GET("/tags", func(c *gin.Context) {
			items, err := services.Tags.List()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, items)
		})

		api.GET("/menus", func(c *gin.Context) {
			items, err := services.Menus.List()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, items)
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

		api.GET("/attachments", func(c *gin.Context) {
			items, err := services.Attachments.List()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, items)
		})

		api.GET("/backups", func(c *gin.Context) {
			items, err := services.Backups.List()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, items)
		})

		api.GET("/settings", func(c *gin.Context) {
			items, err := services.Settings.List()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, items)
		})

		api.GET("/users", func(c *gin.Context) {
			items, err := services.Users.FindAll()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, items)
		})
	}

	staticDir := filepath.Join("web", "dist")
	g.Static("/assets", filepath.Join(staticDir, "assets"))
	g.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
			return
		}
		c.File(filepath.Join(staticDir, "index.html"))
	})
}
