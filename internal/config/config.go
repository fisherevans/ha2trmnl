package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HomeAssistantHost  string `yaml:"home_assistant_host"`
	HomeAssistantToken string `yaml:"home_assistant_token"`
	TrmnlWebhook       string `yaml:"trmnl_webhook"`
	DryRun             bool   `yaml:"dry_run"`
	TimeZone           string `yaml:"time_zone"`
	Debug              bool   `yaml:"debug"`
}

func FromFile(path string) (Config, error) {
	log.Println("Config file: " + path)
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
