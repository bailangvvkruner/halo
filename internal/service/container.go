package service

import (
	"halo/internal/config"

	"gorm.io/gorm"
)

type Container struct {
	Users       *UserService
	Posts       *PostService
	Pages       *PageService
	Categories  *CategoryService
	Tags        *TagService
	Menus       *MenuService
	Comments    *CommentService
	Themes      *ThemeService
	Plugins     *PluginService
	Attachments *AttachmentService
	Backups     *BackupService
	Settings    *SettingService
	Auth        *AuthService
	Dashboard   *DashboardService
	Setup       *SetupService
}

func NewContainer(db *gorm.DB, cfg config.Config) (*Container, error) {
	users := NewUserService(db)
	posts := NewPostService(db)
	pages := NewPageService(db)
	categories := NewCategoryService(db)
	tags := NewTagService(db)
	menus := NewMenuService(db)
	comments := NewCommentService(db)
	themes := NewThemeService(db)
	plugins := NewPluginService(db)
	attachments := NewAttachmentService(db)
	backups := NewBackupService(db)
	settings := NewSettingService(db)
	auth := NewAuthService(users, cfg)
	dashboard := NewDashboardService(db)
	setup := NewSetupService(db, users, settings)

	return &Container{
		Users:       users,
		Posts:       posts,
		Pages:       pages,
		Categories:  categories,
		Tags:        tags,
		Menus:       menus,
		Comments:    comments,
		Themes:      themes,
		Plugins:     plugins,
		Attachments: attachments,
		Backups:     backups,
		Settings:    settings,
		Auth:        auth,
		Dashboard:   dashboard,
		Setup:       setup,
	}, nil
}
