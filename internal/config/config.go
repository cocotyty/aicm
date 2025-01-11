package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	LLMAPIKey string `json:"LLM_API_KEY"`
	LLMModel  string `json:"LLM_MODEL"`
	LLMAPIURL string `json:"LLM_API_URL"`
}

func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, ".aicm", "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func SetConfig(key, value string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".aicm")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.json")

	// 读取现有配置
	cfg := &Config{}
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, cfg); err != nil {
			return err
		}
	}

	// 更新配置
	switch key {
	case "LLM_API_KEY":
		cfg.LLMAPIKey = value
	case "LLM_MODEL":
		cfg.LLMModel = value
	case "LLM_API_URL":
		cfg.LLMAPIURL = value
	default:
		return fmt.Errorf("invalid config key: %s", key)
	}

	// 保存配置
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
