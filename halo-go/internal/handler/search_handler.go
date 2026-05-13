package handler

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"

	"github.com/halo-dev/halo-go/internal/extension"
	"github.com/halo-dev/halo-go/internal/model"
	"github.com/halo-dev/halo-go/internal/service"
)

type SearchHandler struct {
	postService    service.PostService
	pageService    service.SinglePageService
	categoryService service.CategoryService
	tagService     service.TagService
}

func NewSearchHandler(
	postSvc service.PostService,
	pageSvc service.SinglePageService,
	catSvc service.CategoryService,
	tagSvc service.TagService,
) *SearchHandler {
	return &SearchHandler{
		postService:    postSvc,
		pageService:    pageSvc,
		categoryService: catSvc,
		tagService:     tagSvc,
	}
}

func (h *SearchHandler) Search(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请输入搜索关键词",
			"data":    nil,
		})
		return
	}

	size := 20
	page := 1

	results := make([]SearchResult, 0)

	postResult, _ := h.postService.List(c.Request.Context(), &extension.ListOptions{Size: 0})
	for _, item := range postResult.Items {
		if post, ok := item.(*model.Post); ok && h.matchKeyword(post, keyword) {
			year := post.Metadata.CreationTimestamp.Year()
			month := int(post.Metadata.CreationTimestamp.Month())
			day := post.Metadata.CreationTimestamp.Day()
			results = append(results, SearchResult{
				Type:      "Post",
				Name:      post.Metadata.Name,
				Title:     post.Spec.Title,
				Excerpt:   h.getExcerpt(post),
				Permalink: post.Status.Permalink,
				Date: carbon.CreateFromDateTime(year, month, day, 0, 0, 0, "Asia/Shanghai").ToDateTimeString(),
			})
		}
	}

	pageResult, _ := h.pageService.List(c.Request.Context(), &extension.ListOptions{Size: 0})
	for _, item := range pageResult.Items {
		if page, ok := item.(*model.SinglePage); ok && h.matchKeywordPage(page, keyword) {
			year := page.Metadata.CreationTimestamp.Year()
			month := int(page.Metadata.CreationTimestamp.Month())
			day := page.Metadata.CreationTimestamp.Day()
			results = append(results, SearchResult{
				Type:      "SinglePage",
				Name:      page.Metadata.Name,
				Title:     page.Spec.Title,
				Excerpt:   h.getPageExcerpt(page),
				Permalink: page.Status.Permalink,
				Date: carbon.CreateFromDateTime(year, month, day, 0, 0, 0, "Asia/Shanghai").ToDateTimeString(),
			})
		}
	}

	total := len(results)
	start := (page - 1) * size
	end := start + size
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data": gin.H{
			"total": total,
			"page":  page,
			"size":  size,
			"items": results[start:end],
		},
	})
}

type SearchResult struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Excerpt   string `json:"excerpt"`
	Permalink string `json:"permalink"`
	Date      string `json:"date"`
}

func (h *SearchHandler) matchKeyword(post *model.Post, keyword string) bool {
	lowerContent := strings.ToLower(post.Spec.Title + " " + post.Spec.Slug)
	lowerKeyword := strings.ToLower(keyword)
	return strings.Contains(lowerContent, lowerKeyword)
}

func (h *SearchHandler) matchKeywordPage(page *model.SinglePage, keyword string) bool {
	lowerContent := strings.ToLower(page.Spec.Title + " " + page.Spec.Slug)
	lowerKeyword := strings.ToLower(keyword)
	return strings.Contains(lowerContent, lowerKeyword)
}

func (h *SearchHandler) getExcerpt(post *model.Post) string {
	if post.Spec.Excerpt.Raw != "" {
		return truncateText(post.Spec.Excerpt.Raw, 200)
	}
	content := stripHTMLTags(post.Spec.Content)
	return truncateText(content, 200)
}

func (h *SearchHandler) getPageExcerpt(page *model.SinglePage) string {
	if page.Spec.Excerpt.Raw != "" {
		return truncateText(page.Spec.Excerpt.Raw, 200)
	}
	content := stripHTMLTags(page.Spec.Content)
	return truncateText(content, 200)
}

var htmlTagRegex = regexp.MustCompile(`<[^>]*>`)

func stripHTMLTags(s string) string {
	return htmlTagRegex.ReplaceAllString(s, "")
}

func truncateText(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}
