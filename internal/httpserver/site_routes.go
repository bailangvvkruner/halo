package httpserver

import (
	"html/template"
	"net/http"

	"halo/internal/service"

	"github.com/gin-gonic/gin"
)

var siteTemplate = template.Must(template.New("site").Parse(`<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{ .Title }}</title>
  <style>
    body { font-family: Inter, system-ui, sans-serif; margin: 0; background: #f8fafc; color: #0f172a; }
    .container { max-width: 880px; margin: 0 auto; padding: 48px 24px 72px; }
    .hero { margin-bottom: 32px; }
    .hero h1 { margin: 0 0 10px; font-size: 36px; }
    .meta { color: #64748b; font-size: 14px; }
    .content { background: #fff; border: 1px solid #e2e8f0; border-radius: 16px; padding: 28px; line-height: 1.8; }
    .back { display: inline-block; margin-top: 24px; color: #2563eb; text-decoration: none; }
    .list { display: grid; gap: 16px; }
    .card { background: #fff; border: 1px solid #e2e8f0; border-radius: 14px; padding: 20px; }
    .card h2 { margin: 0 0 8px; }
    .excerpt { color: #475569; }
  </style>
</head>
<body>
  <main class="container">
    {{ if .List }}
      <section class="hero">
        <h1>{{ .Title }}</h1>
        <p class="meta">基础公开渲染入口</p>
      </section>
      <section class="list">
        {{ range .Items }}
          <article class="card">
            <h2>{{ .Title }}</h2>
            <p class="excerpt">{{ .Excerpt }}</p>
            <a class="back" href="{{ .URL }}">查看详情</a>
          </article>
        {{ end }}
      </section>
    {{ else }}
      <section class="hero">
        <h1>{{ .Title }}</h1>
        <p class="meta">{{ .Meta }}</p>
      </section>
      <article class="content">{{ .Content }}</article>
      <a class="back" href="/">返回首页</a>
    {{ end }}
  </main>
</body>
</html>`))

func registerSiteRoutes(g *gin.Engine, services *service.Container) {
	g.GET("/posts/:slug", func(c *gin.Context) {
		post, err := services.Posts.FindBySlug(c.Param("slug"))
		if err != nil {
			c.String(http.StatusNotFound, "post not found")
			return
		}
		c.Status(http.StatusOK)
		siteTemplate.Execute(c.Writer, gin.H{
			"Title":   post.Title,
			"Meta":    post.Category,
			"Content": template.HTML(post.Content),
		})
	})

	g.GET("/pages/:slug", func(c *gin.Context) {
		pages, err := services.Pages.List()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		for _, page := range pages {
			if page.Slug == c.Param("slug") {
				c.Status(http.StatusOK)
				siteTemplate.Execute(c.Writer, gin.H{
					"Title":   page.Title,
					"Meta":    page.Slug,
					"Content": template.HTML(page.Content),
				})
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
		items := make([]gin.H, 0, len(posts))
		for _, post := range posts {
			items = append(items, gin.H{
				"Title":   post.Title,
				"Excerpt": post.Excerpt,
				"URL":     "/posts/" + post.Slug,
			})
		}
		c.Status(http.StatusOK)
		siteTemplate.Execute(c.Writer, gin.H{
			"Title": "Halo Go",
			"List":  true,
			"Items": items,
		})
	})
}
