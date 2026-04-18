package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"halo/internal/config"
)

type ThemeRenderService struct {
	cfg config.Config
	themes *ThemeService
}

func NewThemeRenderService(cfg config.Config, themes *ThemeService) *ThemeRenderService {
	return &ThemeRenderService{cfg: cfg, themes: themes}
}

func (s *ThemeRenderService) ActiveThemeName() (string, error) {
	items, err := s.themes.List()
	if err != nil {
		return "", err
	}
	for _, item := range items {
		if item.Activated {
			return item.Name, nil
		}
	}
	return "default", nil
}

func (s *ThemeRenderService) EnsureDefaultThemeFiles() error {
	themeDir := filepath.Join(s.cfg.WorkDir, "themes", "default")
	if err := os.MkdirAll(themeDir, 0o755); err != nil {
		return err
	}

	files := map[string]string{
		"index.html": defaultIndexTemplate,
		"post.html": defaultPostTemplate,
		"page.html": defaultPageTemplate,
	}

	for name, content := range files {
		path := filepath.Join(themeDir, name)
		if _, err := os.Stat(path); err == nil {
			continue
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return err
		}
	}

	return nil
}

func (s *ThemeRenderService) LoadTemplate(name string) (string, error) {
	themeName, err := s.ActiveThemeName()
	if err != nil {
		return "", err
	}
	path := filepath.Join(s.cfg.WorkDir, "themes", themeName, name)
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("load theme template: %w", err)
	}
	return string(content), nil
}

func (s *ThemeRenderService) Render(templateContent string, data map[string]string) string {
	result := templateContent
	for key, value := range data {
		result = strings.ReplaceAll(result, "{{ "+key+" }}", value)
		result = strings.ReplaceAll(result, "{{"+key+"}}", value)
	}
	return result
}

const defaultIndexTemplate = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{ title }}</title>
  <style>
    body { font-family: Inter, system-ui, sans-serif; margin: 0; background: #f8fafc; color: #0f172a; }
    .container { max-width: 880px; margin: 0 auto; padding: 48px 24px 72px; }
    .card { background: #fff; border: 1px solid #e2e8f0; border-radius: 14px; padding: 20px; margin-bottom: 16px; }
  </style>
</head>
<body>
  <main class="container">
    <h1>{{ title }}</h1>
    {{ content }}
  </main>
</body>
</html>`

const defaultPostTemplate = defaultIndexTemplate
const defaultPageTemplate = defaultIndexTemplate
