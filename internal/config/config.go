package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HomeAssistantConfig `yaml:"home_assistant"`
	ServeConfig         `yaml:"serve"`
	PushConfig          `yaml:"push"`

	TimeZone string `yaml:"time_zone"`
	Debug    bool   `yaml:"debug"`
}

type HomeAssistantConfig struct {
	HomeAssistantHost  string `yaml:"host"`
	HomeAssistantToken string `yaml:"token"`
}

type ServeConfig struct {
	ApiToken string `yaml:"api_token"`
	Address  string `yaml:"address"`
}

type PushConfig struct {
	Interval string `yaml:"interval"`
	Webhook  string `yaml:"webhook"`
	DryRun   bool   `yaml:"dry_run"`
}

func FromFile(path string) (Config, error) {
	var c Config

	configContents, err := os.ReadFile(path)
	if err != nil {
		return c, fmt.Errorf("unable to read config file: %w", err)
	}

	if err = yaml.Unmarshal(configContents, &c); err != nil {
		return c, fmt.Errorf("unable to parse config file: %w", err)
	}

	return c, nil
}
