package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Addr        string
	WorkDir     string
	DatabaseURL string
	JWTSecret   string
	SessionKey  string
}

func Load() Config {
	workDir := getenv("HALO_WORK_DIR", filepath.Join(userHomeDir(), ".halo2"))
	databaseURL := getenv("HALO_DB_PATH", filepath.Join(workDir, "db", "halo.db"))

	return Config{
		Addr:        getenv("HALO_ADDR", ":8090"),
		WorkDir:     workDir,
		DatabaseURL: databaseURL,
		JWTSecret:   getenv("HALO_JWT_SECRET", "halo-go-secret"),
		SessionKey:  getenv("HALO_SESSION_KEY", "halo-go-session"),
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func userHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return home
}
