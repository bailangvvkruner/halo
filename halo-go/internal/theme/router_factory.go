package theme

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/halo-dev/halo-go/internal/model"
	"github.com/halo-dev/halo-go/internal/service"
)

type RouteFactory interface {
	CreateRoute() gin.HandlerFunc
	GetPath() string
}

type IndexRouteFactory struct {
	templateEngine TemplateEngine
	postService    service.PostService
}

func NewIndexRouteFactory(templateEngine TemplateEngine, postService service.PostService) *IndexRouteFactory {
	return &IndexRouteFactory{
		templateEngine: templateEngine,
		postService:    postService,
	}
}

func (f *IndexRouteFactory) GetPath() string {
	return "/"
}

func (f *IndexRouteFactory) CreateRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := f.postService.List(c.Request.Context(), nil)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		posts := make([]*model.Post, 0)
		for _, ext := range result.Items {
			post := ext.(*model.Post)
			if post.Spec.Publish && !post.Spec.Deleted {
				posts = append(posts, post)
			}
		}

		data := map[string]interface{}{
			"posts": posts,
		}

		html, err := f.templateEngine.Render("index.html", data)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	}
}

type PostRouteFactory struct {
	templateEngine TemplateEngine
	postService    service.PostService
}

func NewPostRouteFactory(templateEngine TemplateEngine, postService service.PostService) *PostRouteFactory {
	return &PostRouteFactory{
		templateEngine: templateEngine,
		postService:    postService,
	}
}

func (f *PostRouteFactory) GetPath() string {
	return "/archives/:slug"
}

func (f *PostRouteFactory) CreateRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")
		result, err := f.postService.List(c.Request.Context(), nil)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		var post *model.Post
		for _, ext := range result.Items {
			p := ext.(*model.Post)
			if p.Spec.Slug == slug && p.Spec.Publish && !p.Spec.Deleted {
				post = p
				break
			}
		}

		if post == nil {
			c.Status(http.StatusNotFound)
			return
		}

		data := map[string]interface{}{
			"post": post,
		}

		html, err := f.templateEngine.Render("post.html", data)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	}
}

type ArchiveRouteFactory struct {
	templateEngine TemplateEngine
	postService    service.PostService
}

func NewArchiveRouteFactory(templateEngine TemplateEngine, postService service.PostService) *ArchiveRouteFactory {
	return &ArchiveRouteFactory{
		templateEngine: templateEngine,
		postService:    postService,
	}
}

func (f *ArchiveRouteFactory) GetPath() string {
	return "/archives"
}

func (f *ArchiveRouteFactory) CreateRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := f.postService.List(c.Request.Context(), nil)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		posts := make([]*model.Post, 0)
		for _, ext := range result.Items {
			post := ext.(*model.Post)
			if post.Spec.Publish && !post.Spec.Deleted {
				posts = append(posts, post)
			}
		}

		data := map[string]interface{}{
			"posts": posts,
		}

		html, err := f.templateEngine.Render("archives.html", data)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	}
}

type TagsRouteFactory struct {
	templateEngine TemplateEngine
	tagService     service.TagService
}

func NewTagsRouteFactory(templateEngine TemplateEngine, tagService service.TagService) *TagsRouteFactory {
	return &TagsRouteFactory{
		templateEngine: templateEngine,
		tagService:     tagService,
	}
}

func (f *TagsRouteFactory) GetPath() string {
	return "/tags"
}

func (f *TagsRouteFactory) CreateRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := f.tagService.List(c.Request.Context(), nil)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		tags := make([]*model.Tag, 0)
		for _, ext := range result.Items {
			tags = append(tags, ext.(*model.Tag))
		}

		data := map[string]interface{}{
			"tags": tags,
		}

		html, err := f.templateEngine.Render("tags.html", data)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	}
}

type TagPostRouteFactory struct {
	templateEngine TemplateEngine
	tagService     service.TagService
	postService    service.PostService
}

func NewTagPostRouteFactory(templateEngine TemplateEngine, tagService service.TagService, postService service.PostService) *TagPostRouteFactory {
	return &TagPostRouteFactory{
		templateEngine: templateEngine,
		tagService:     tagService,
		postService:    postService,
	}
}

func (f *TagPostRouteFactory) GetPath() string {
	return "/tags/:slug"
}

func (f *TagPostRouteFactory) CreateRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")
		tagResult, err := f.tagService.List(c.Request.Context(), nil)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		var tag *model.Tag
		for _, ext := range tagResult.Items {
			t := ext.(*model.Tag)
			if t.Spec.Slug == slug {
				tag = t
				break
			}
		}

		if tag == nil {
			c.Status(http.StatusNotFound)
			return
		}

		postResult, err := f.postService.List(c.Request.Context(), nil)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		posts := make([]*model.Post, 0)
		for _, ext := range postResult.Items {
			post := ext.(*model.Post)
			if post.Spec.Publish && !post.Spec.Deleted {
				posts = append(posts, post)
			}
		}

		data := map[string]interface{}{
			"tag":   tag,
			"posts": posts,
		}

		html, err := f.templateEngine.Render("tag.html", data)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	}
}

type CategoriesRouteFactory struct {
	templateEngine   TemplateEngine
	categoryService  service.CategoryService
}

func NewCategoriesRouteFactory(templateEngine TemplateEngine, categoryService service.CategoryService) *CategoriesRouteFactory {
	return &CategoriesRouteFactory{
		templateEngine:  templateEngine,
		categoryService: categoryService,
	}
}

func (f *CategoriesRouteFactory) GetPath() string {
	return "/categories"
}

func (f *CategoriesRouteFactory) CreateRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := f.categoryService.List(c.Request.Context(), nil)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		categories := make([]*model.Category, 0)
		for _, ext := range result.Items {
			categories = append(categories, ext.(*model.Category))
		}

		data := map[string]interface{}{
			"categories": categories,
		}

		html, err := f.templateEngine.Render("categories.html", data)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	}
}

type CategoryPostRouteFactory struct {
	templateEngine   TemplateEngine
	categoryService  service.CategoryService
	postService      service.PostService
}

func NewCategoryPostRouteFactory(templateEngine TemplateEngine, categoryService service.CategoryService, postService service.PostService) *CategoryPostRouteFactory {
	return &CategoryPostRouteFactory{
		templateEngine:  templateEngine,
		categoryService: categoryService,
		postService:     postService,
	}
}

func (f *CategoryPostRouteFactory) GetPath() string {
	return "/categories/:slug"
}

func (f *CategoryPostRouteFactory) CreateRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")
		catResult, err := f.categoryService.List(c.Request.Context(), nil)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		var category *model.Category
		for _, ext := range catResult.Items {
			cat := ext.(*model.Category)
			if cat.Spec.Slug == slug {
				category = cat
				break
			}
		}

		if category == nil {
			c.Status(http.StatusNotFound)
			return
		}

		postResult, err := f.postService.List(c.Request.Context(), nil)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		posts := make([]*model.Post, 0)
		for _, ext := range postResult.Items {
			post := ext.(*model.Post)
			if post.Spec.Publish && !post.Spec.Deleted {
				posts = append(posts, post)
			}
		}

		data := map[string]interface{}{
			"category": category,
			"posts":    posts,
		}

		html, err := f.templateEngine.Render("category.html", data)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	}
}

type SinglePageRouteFactory struct {
	templateEngine   TemplateEngine
	singlePageService service.SinglePageService
}

func NewSinglePageRouteFactory(templateEngine TemplateEngine, singlePageService service.SinglePageService) *SinglePageRouteFactory {
	return &SinglePageRouteFactory{
		templateEngine:    templateEngine,
		singlePageService: singlePageService,
	}
}

func (f *SinglePageRouteFactory) GetPath() string {
	return "/:slug"
}

func (f *SinglePageRouteFactory) CreateRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")
		result, err := f.singlePageService.List(c.Request.Context(), nil)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		var page *model.SinglePage
		for _, ext := range result.Items {
			p := ext.(*model.SinglePage)
			if p.Spec.Slug == slug && p.Spec.Publish && !p.Spec.Deleted {
				page = p
				break
			}
		}

		if page == nil {
			c.Status(http.StatusNotFound)
			return
		}

		data := map[string]interface{}{
			"page": page,
		}

		html, err := f.templateEngine.Render("page.html", data)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	}
}