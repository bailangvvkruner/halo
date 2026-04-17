package httpserver

import "halo/internal/model"

func toPost(title, slug, content, excerpt string, published bool) *model.Post {
	return &model.Post{
		Title:     title,
		Slug:      slug,
		Content:   content,
		Excerpt:   excerpt,
		Published: published,
	}
}
