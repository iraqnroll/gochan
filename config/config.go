package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

const (
	CONFIG_FILENAME = "config/config.toml"
)

type (
	Config struct {
		Global   Global
		Database PostgresConfig
		Api      API
		Frontend Frontend
	}

	Global struct {
		Shortname         string
		Subtitle          string
		RecentPostsNum    int
		AllowedMediaTypes []string
	}

	Frontend struct {
		StaticDir string
		Enabled   bool
	}

	API struct {
		Enabled          bool
		RecentPostsNum   int
		SessionTokenSize int
	}
)

func (cfg Global) String() string {
	return fmt.Sprintf("shortname=%s subtitle=%s", cfg.Shortname, cfg.Subtitle)
}

func InitConfig() *Config {
	var config Config
	_, err := toml.DecodeFile(CONFIG_FILENAME, &config)
	if err != nil {
		fmt.Println("Error decoding configuration file :", err)
		os.Exit(1)
	}

	return &config
}
