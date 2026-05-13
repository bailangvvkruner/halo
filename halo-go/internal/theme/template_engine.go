package theme

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"path"
	"sync"
	"time"

	"github.com/gomarkdown/markdown"
)

type TemplateEngine interface {
	Render(templateName string, data map[string]interface{}) (string, error)
	LoadTemplates(themeFS fs.FS) error
}

type templateEngine struct {
	templates *template.Template
	cache     sync.Map
}

func NewTemplateEngine() TemplateEngine {
	return &templateEngine{}
}

func (e *templateEngine) LoadTemplates(themeFS fs.FS) error {
	tmpl := template.New("").Funcs(template.FuncMap{
		"date": func(t *time.Time, format string) string {
			if t == nil {
				return ""
			}
			return t.Format(format)
		},
		"truncate": func(s string, length int) string {
			if len(s) <= length {
				return s
			}
			return s[:length] + "..."
		},
		"markdown": func(s string) template.HTML {
			return template.HTML(markdown.ToHTML([]byte(s), nil, nil))
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"safeJS": func(s string) template.JS {
			return template.JS(s)
		},
		"urlJoin": func(parts ...string) string {
			return path.Join(parts...)
		},
		"eq": func(a, b interface{}) bool {
			return a == b
		},
		"ne": func(a, b interface{}) bool {
			return a != b
		},
		"gt": func(a, b int) bool {
			return a > b
		},
		"lt": func(a, b int) bool {
			return a < b
		},
	})

	err := fs.WalkDir(themeFS, ".", func(filepath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && path.Ext(filepath) == ".html" {
			content, err := fs.ReadFile(themeFS, filepath)
			if err != nil {
				return fmt.Errorf("failed to read template %s: %w", filepath, err)
			}
			tmpl, err = tmpl.New(filepath).Parse(string(content))
			if err != nil {
				return fmt.Errorf("failed to parse template %s: %w", filepath, err)
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	e.templates = tmpl
	return nil
}

func (e *templateEngine) Render(templateName string, data map[string]interface{}) (string, error) {
	if e.templates == nil {
		return "", fmt.Errorf("templates not loaded")
	}

	tmpl := e.templates.Lookup(templateName)
	if tmpl == nil {
		return "", fmt.Errorf("template %s not found", templateName)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return buf.String(), nil
}