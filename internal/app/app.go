package app

import (
	"fmt"
	"os"
	"path/filepath"

	"halo/internal/config"
	"halo/internal/database"
	"halo/internal/httpserver"
	"halo/internal/seed"
	"halo/internal/service"
)

type App struct {
	server *httpserver.Server
	config config.Config
}

func New() (*App, error) {
	cfg := config.Load()
	if err := ensureWorkDirs(cfg.WorkDir); err != nil {
		return nil, err
	}

	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if err := seed.Bootstrap(db); err != nil {
		return nil, err
	}

	services, err := service.NewContainer(db, cfg)
	if err != nil {
		return nil, err
	}

	server := httpserver.New(cfg, services)

	return &App{server: server, config: cfg}, nil
}

func (a *App) Run() error {
	return a.server.Run(a.config.Addr)
}

func ensureWorkDirs(root string) error {
	dirs := []string{
		root,
		filepath.Join(root, "db"),
		filepath.Join(root, "logs"),
		filepath.Join(root, "themes"),
		filepath.Join(root, "plugins"),
		filepath.Join(root, "attachments"),
		filepath.Join(root, "backups"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create work dir %s: %w", dir, err)
		}
	}

	return nil
}
