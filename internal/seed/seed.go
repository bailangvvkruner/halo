package seed

import (
	"halo/internal/model"

	"gorm.io/gorm"
)

func Bootstrap(db *gorm.DB) error {
	if err := ensureTheme(db); err != nil {
		return err
	}
	if err := ensurePlugin(db); err != nil {
		return err
	}
	if err := ensureSettings(db); err != nil {
		return err
	}
	if err := ensureSamplePost(db); err != nil {
		return err
	}
	if err := ensurePage(db); err != nil {
		return err
	}
	if err := ensureCategory(db); err != nil {
		return err
	}
	if err := ensureTag(db); err != nil {
		return err
	}
	if err := ensureMenu(db); err != nil {
		return err
	}
	if err := ensureBackup(db); err != nil {
		return err
	}
	return nil
}

func ensureTheme(db *gorm.DB) error {
	return firstOrCreate(db, &model.Theme{}, model.Theme{
		Name:        "default",
		DisplayName: "Default Theme",
		Activated:   true,
	})
}

func ensurePlugin(db *gorm.DB) error {
	return firstOrCreate(db, &model.Plugin{}, model.Plugin{
		Name:        "content-core",
		DisplayName: "Content Core",
		Enabled:     true,
	})
}

func ensureSettings(db *gorm.DB) error {
	if err := firstOrCreate(db, &model.Setting{}, model.Setting{Key: "site.title", Value: "Halo Go"}); err != nil {
		return err
	}
	return firstOrCreate(db, &model.Setting{}, model.Setting{Key: "site.subtitle", Value: "Modern Go CMS"})
}

func ensureSamplePost(db *gorm.DB) error {
	return firstOrCreate(db, &model.Post{}, model.Post{
		Title:     "欢迎使用 Halo Go",
		Slug:      "welcome-to-halo-go",
		Excerpt:   "这是 Go 重构版的初始文章，用于验证内容链路和前后端联通。",
		Content:   "# Halo Go\n\n当前版本已经完成旧项目归档、新后端骨架、新前端骨架以及 Docker 支持。",
		Template:  "post",
		Category:  "news",
		Tags:      "go,halo,cms",
		Published: true,
	})
}

func ensurePage(db *gorm.DB) error {
	return firstOrCreate(db, &model.Page{}, model.Page{
		Title:     "关于我们",
		Slug:      "about",
		Content:   "这是 Go 重构版的示例单页。",
		Published: true,
	})
}

func ensureCategory(db *gorm.DB) error {
	return firstOrCreate(db, &model.Category{}, model.Category{
		Name:        "news",
		Slug:        "news",
		DisplayName: "项目动态",
	})
}

func ensureTag(db *gorm.DB) error {
	return firstOrCreate(db, &model.Tag{}, model.Tag{
		Name: "go",
		Slug: "go",
	})
}

func ensureMenu(db *gorm.DB) error {
	return firstOrCreate(db, &model.Menu{}, model.Menu{
		Name:  "main",
		Items: `[{"label":"首页","path":"/"},{"label":"关于","path":"/about"}]`,
	})
}

func ensureBackup(db *gorm.DB) error {
	return firstOrCreate(db, &model.Backup{}, model.Backup{
		Filename: "initial-backup.zip",
		Status:   "Succeeded",
		Size:     0,
	})
}

func firstOrCreate[T any](db *gorm.DB, query any, value T) error {
	return db.Where(query).FirstOrCreate(&value).Error
}
