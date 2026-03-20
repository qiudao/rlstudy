package rlstudy

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port   int    `json:"port"`
	Arms   int    `json:"arms"`
	Socket string `json:"socket"`
}

func DefaultConfig() Config {
	return Config{Port: 21320, Arms: 10, Socket: "/tmp/bandit-env.sock"}
}

func LoadConfig(path string) (Config, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}
