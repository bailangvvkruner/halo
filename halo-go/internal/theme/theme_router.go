package theme

import (
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/halo-dev/halo-go/internal/service"
)

type ThemeRouteRules struct {
	EnabledRoutes []string
}

type ThemeCompositeRouter struct {
	templateEngine     TemplateEngine
	themeFinder        ThemeFinder
	postService        service.PostService
	tagService         service.TagService
	categoryService    service.CategoryService
	singlePageService  service.SinglePageService
	workDir            string
}

func NewThemeCompositeRouter(
	templateEngine TemplateEngine,
	themeFinder ThemeFinder,
	postService service.PostService,
	tagService service.TagService,
	categoryService service.CategoryService,
	singlePageService service.SinglePageService,
	workDir string,
) *ThemeCompositeRouter {
	return &ThemeCompositeRouter{
		templateEngine:    templateEngine,
		themeFinder:       themeFinder,
		postService:       postService,
		tagService:        tagService,
		categoryService:   categoryService,
		singlePageService: singlePageService,
		workDir:           workDir,
	}
}

func (r *ThemeCompositeRouter) RegisterRoutes(engine *gin.Engine) error {
	theme, err := r.themeFinder.GetActiveTheme(nil)
	if err != nil {
		return err
	}

	themePath := path.Join(r.workDir, "themes", theme.Metadata.Name)
	_, err = os.Stat(themePath)
	if err != nil {
		return err
	}

	filesystem := os.DirFS(themePath)
	err = r.templateEngine.LoadTemplates(filesystem)
	if err != nil {
		return err
	}

	factories := []RouteFactory{
		NewIndexRouteFactory(r.templateEngine, r.postService),
		NewArchiveRouteFactory(r.templateEngine, r.postService),
		NewPostRouteFactory(r.templateEngine, r.postService),
		NewTagsRouteFactory(r.templateEngine, r.tagService),
		NewTagPostRouteFactory(r.templateEngine, r.tagService, r.postService),
		NewCategoriesRouteFactory(r.templateEngine, r.categoryService),
		NewCategoryPostRouteFactory(r.templateEngine, r.categoryService, r.postService),
		NewSinglePageRouteFactory(r.templateEngine, r.singlePageService),
	}

	for _, factory := range factories {
		engine.GET(factory.GetPath(), factory.CreateRoute())
	}

	engine.GET("/themes/:themeName/assets/*filepath", func(c *gin.Context) {
		themeName := c.Param("themeName")
		filepath := c.Param("filepath")
		assetPath := path.Join(r.workDir, "themes", themeName, "assets", filepath)

		file, err := os.Open(assetPath)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		contentType := r.guessContentType(filepath)
		c.Header("Content-Type", contentType)
		c.Header("Cache-Control", "public, max-age=31536000")
		http.ServeContent(c.Writer, c.Request, fileInfo.Name(), fileInfo.ModTime(), file)
	})

	return nil
}

func (r *ThemeCompositeRouter) guessContentType(name string) string {
	ext := path.Ext(name)
	switch ext {
	case ".ico":
		return "image/x-icon"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".svg":
		return "image/svg+xml"
	case ".css":
		return "text/css; charset=utf-8"
	case ".js", ".mjs":
		return "application/javascript"
	case ".woff2":
		return "font/woff2"
	case ".woff":
		return "font/woff"
	case ".ttf":
		return "font/ttf"
	case ".eot":
		return "application/vnd.ms-fontobject"
	case ".json":
		return "application/json"
	case ".html", ".htm":
		return "text/html; charset=utf-8"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".ogg":
		return "audio/ogg"
	case ".mp3":
		return "audio/mpeg"
	default:
		return "application/octet-stream"
	}
}