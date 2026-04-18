package httpserver

import (
	"net/http"
	"strings"

	"halo/internal/service"

	"github.com/gin-gonic/gin"
)

func registerSiteRoutes(g *gin.Engine, services *service.Container) {
	g.GET("/posts/:slug", func(c *gin.Context) {
		post, err := services.Posts.FindBySlug(c.Param("slug"))
		if err != nil {
			c.String(http.StatusNotFound, "post not found")
			return
		}
		tpl, err := services.ThemeRender.LoadTemplate("post.html")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		content := services.ThemeRender.Render(tpl, map[string]string{
			"title":   post.Title,
			"content": post.Content,
			"meta":    post.Category,
			"subtitle": post.Excerpt,
		})
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
	})

	g.GET("/pages/:slug", func(c *gin.Context) {
		pages, err := services.Pages.List()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		for _, page := range pages {
			if page.Slug == c.Param("slug") {
				tpl, err := services.ThemeRender.LoadTemplate("page.html")
				if err != nil {
					c.String(http.StatusInternalServerError, err.Error())
					return
				}
				content := services.ThemeRender.Render(tpl, map[string]string{
					"title":   page.Title,
					"content": page.Content,
					"meta":    page.Slug,
					"subtitle": page.Slug,
				})
				c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
				return
			}
		}
		c.String(http.StatusNotFound, "page not found")
	})

	g.GET("/", func(c *gin.Context) {
		posts, err := services.Posts.List()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		items := make([]string, 0, len(posts))
		for _, post := range posts {
			items = append(items, "<article class=\"card\"><h2>"+post.Title+"</h2><p class=\"excerpt\">"+post.Excerpt+"</p><a class=\"back\" href=\"/posts/"+post.Slug+"\">查看详情</a></article>")
		}
		tpl, err := services.ThemeRender.LoadTemplate("index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		content := services.ThemeRender.Render(tpl, map[string]string{
			"title":   "Halo Go",
			"subtitle": "基于 Go 的公开渲染入口",
			"content": strings.Join(items, ""),
		})
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
	})
}
