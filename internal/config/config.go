package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	ListenOn   string
	SqlitePath string
	DrawCount  int64 // 凑够多少个人就开奖
}

func MustLoad(filepath string) Config {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}
	cfg := Config{}
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		log.Fatalf("failed to parse config file: %v", err)
	}
	if cfg.ListenOn == "" {
		log.Fatalf("listen_on is empty")
	}
	if cfg.SqlitePath == "" {
		log.Fatalf("sqlite_path is empty")
	}
	if cfg.DrawCount <= 0 {
		log.Fatalf("draw_count is empty")
	}
	return cfg
}
