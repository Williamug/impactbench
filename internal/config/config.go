package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	BaseURL    string             `mapstructure:"base_url"`
	Framework  string             `mapstructure:"framework"`
	Storage    StorageConfig      `mapstructure:"storage"`
	Thresholds ThresholdsConfig   `mapstructure:"thresholds"`
	LoadTest   LoadTestConfig     `mapstructure:"load_test"`
}

type StorageConfig struct {
	Type string `mapstructure:"type"`
	Path string `mapstructure:"path"`
}

type ThresholdsConfig struct {
	ResponseTime float64 `mapstructure:"response_time"`
	ErrorRate    float64 `mapstructure:"error_rate"`
}

type LoadTestConfig struct {
	DefaultUsers    int `mapstructure:"default_users"`
	DefaultDuration int `mapstructure:"default_duration"`
}

func Load() (*Config, error) {
	v := viper.New()

	// Set Defaults
	v.SetDefault("framework", "http")
	v.SetDefault("storage.type", "sqlite")
	v.SetDefault("storage.path", "./impactbench.db")
	v.SetDefault("thresholds.response_time", 10.0) // 10%
	v.SetDefault("thresholds.error_rate", 1.0)      // 1%
	v.SetDefault("load_test.default_users", 10)
	v.SetDefault("load_test.default_duration", 30)

	// Global Config: ~/.impactbench/config.yaml
	home, err := os.UserHomeDir()
	if err == nil {
		globalConfig := filepath.Join(home, ".impactbench", "config.yaml")
		v.SetConfigFile(globalConfig)
		if err := v.MergeInConfig(); err == nil {
			// Successfully merged global config
		}
	}

	// Project Config: ./.impactbench/config.yaml
	projectConfig := filepath.Join(".", ".impactbench", "config.yaml")
	v.SetConfigFile(projectConfig)
	if err := v.MergeInConfig(); err == nil {
		// Successfully merged project config (overrides global)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return &cfg, nil
}
