package service

import "gorm.io/gorm"

type DashboardStats struct {
	Posts      int64 `json:"posts"`
	Pages      int64 `json:"pages"`
	Categories int64 `json:"categories"`
	Tags       int64 `json:"tags"`
	Users      int64 `json:"users"`
	Themes     int64 `json:"themes"`
	Plugins    int64 `json:"plugins"`
}

type DashboardService struct {
	db *gorm.DB
}

func NewDashboardService(db *gorm.DB) *DashboardService {
	return &DashboardService{db: db}
}

func (s *DashboardService) Stats(container *Container) (DashboardStats, error) {
	var stats DashboardStats

	if posts, err := container.Posts.List(); err == nil {
		stats.Posts = int64(len(posts))
	} else {
		return stats, err
	}

	if pages, err := container.Pages.List(); err == nil {
		stats.Pages = int64(len(pages))
	} else {
		return stats, err
	}

	if categories, err := container.Categories.List(); err == nil {
		stats.Categories = int64(len(categories))
	} else {
		return stats, err
	}

	if tags, err := container.Tags.List(); err == nil {
		stats.Tags = int64(len(tags))
	} else {
		return stats, err
	}

	if users, err := container.Users.FindAll(); err == nil {
		stats.Users = int64(len(users))
	} else {
		return stats, err
	}

	if themes, err := container.Themes.List(); err == nil {
		stats.Themes = int64(len(themes))
	} else {
		return stats, err
	}

	if plugins, err := container.Plugins.List(); err == nil {
		stats.Plugins = int64(len(plugins))
	} else {
		return stats, err
	}

	return stats, nil
}
