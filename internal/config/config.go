package config

import (
	"os"
	"time"
)

type Config struct {
	ListenAddr, DatabasePath, SessionSecret string
	SessionTTL, WorkerInterval              time.Duration
}

func Load() Config {
	return Config{ListenAddr: env("LISTEN_ADDR", ":8080"), DatabasePath: env("DATABASE_PATH", "./data/dataforge.db"), SessionSecret: env("SESSION_SECRET", "dev-secret"), SessionTTL: duration("SESSION_TTL", 8*time.Hour), WorkerInterval: duration("WORKER_INTERVAL", 2*time.Second)}
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func duration(k string, d time.Duration) time.Duration {
	if v := os.Getenv(k); v != "" {
		if x, e := time.ParseDuration(v); e == nil {
			return x
		}
	}
	return d
}
