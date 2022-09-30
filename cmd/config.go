package main

import (
	"os"
)

func loadConfig(cfg *Config) *Config {
	cfg.Version = os.Getenv("version")
	cfg.App = os.Getenv("app")
	cfg.Server.Port = os.Getenv("port")
	cfg.Server.Env = os.Getenv("env")
	cfg.DB.DSN = os.Getenv("dsn")
	return cfg
}
