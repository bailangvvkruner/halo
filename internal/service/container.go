package service

import (
	"halo/internal/config"

	"gorm.io/gorm"
)

type Container struct {
	Users        *UserService
	Posts        *PostService
	Pages        *PageService
	Categories   *CategoryService
	Tags         *TagService
	Menus        *MenuService
	Comments     *CommentService
	Themes       *ThemeService
	Plugins      *PluginService
	Attachments  *AttachmentService
	Backups      *BackupService
	Settings     *SettingService
	Auth         *AuthService
	Dashboard    *DashboardService
	Setup        *SetupService
	Replies      *ReplyService
	Registration *RegistrationService
	Roles        *RoleService
	ThemeRender  *ThemeRenderService
}

func NewContainer(db *gorm.DB, cfg config.Config) (*Container, error) {
	users := NewUserService(db)
	posts := NewPostService(db)
	pages := NewPageService(db)
	categories := NewCategoryService(db)
	tags := NewTagService(db)
	menus := NewMenuService(db)
	comments := NewCommentService(db)
	themes := NewThemeService(db, cfg)
	plugins := NewPluginService(db, cfg)
	attachments := NewAttachmentService(db)
	backups := NewBackupService(db)
	settings := NewSettingService(db)
	auth := NewAuthService(users, cfg)
	dashboard := NewDashboardService(db)
	setup := NewSetupService(db, users, settings)
	replies := NewReplyService(db)
	registration := NewRegistrationService(db)
	roles := NewRoleService(db)
	themeRender := NewThemeRenderService(cfg, themes)
	if err := roles.EnsureDefaults(); err != nil {
		return nil, err
	}
	if err := themeRender.EnsureDefaultThemeFiles(); err != nil {
		return nil, err
	}

	return &Container{
		Users:        users,
		Posts:        posts,
		Pages:        pages,
		Categories:   categories,
		Tags:         tags,
		Menus:        menus,
		Comments:     comments,
		Themes:       themes,
		Plugins:      plugins,
		Attachments:  attachments,
		Backups:      backups,
		Settings:     settings,
		Auth:         auth,
		Dashboard:    dashboard,
		Setup:        setup,
		Replies:      replies,
		Registration: registration,
		Roles:        roles,
		ThemeRender:  themeRender,
	}, nil
}
