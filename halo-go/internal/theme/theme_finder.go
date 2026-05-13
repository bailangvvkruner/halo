package theme

import (
	"context"
	"fmt"

	"github.com/halo-dev/halo-go/internal/model"
	"github.com/halo-dev/halo-go/internal/service"
)

type ThemeFinder interface {
	GetActiveTheme(ctx context.Context) (*model.Theme, error)
	ListThemes(ctx context.Context) ([]*model.Theme, error)
}

type themeFinder struct {
	themeService service.ThemeService
}

func NewThemeFinder(themeService service.ThemeService) ThemeFinder {
	return &themeFinder{themeService: themeService}
}

func (f *themeFinder) GetActiveTheme(ctx context.Context) (*model.Theme, error) {
	result, err := f.themeService.List(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list themes: %w", err)
	}

	for _, ext := range result.Items {
		theme, ok := ext.(*model.Theme)
		if !ok {
			continue
		}
		if theme.Spec.Active {
			return theme, nil
		}
	}

	return nil, fmt.Errorf("no active theme found")
}

func (f *themeFinder) ListThemes(ctx context.Context) ([]*model.Theme, error) {
	result, err := f.themeService.List(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list themes: %w", err)
	}

	themes := make([]*model.Theme, 0, len(result.Items))
	for _, ext := range result.Items {
		theme, ok := ext.(*model.Theme)
		if !ok {
			continue
		}
		themes = append(themes, theme)
	}

	return themes, nil
}