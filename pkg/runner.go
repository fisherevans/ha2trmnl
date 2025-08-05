package pkg

import (
	"fmt"
)

type Instance struct {
	HomeAssistantHost  string `yaml:"home_assistant_host"`
	HomeAssistantToken string `yaml:"home_assistant_token"`
	TrmnlWebhook       string `yaml:"trmnl_webhook"`
	DryRun             bool   `yaml:"dry_run"`
	Debug              bool   `yaml:"debug"`
}

func (i Instance) Run() error {
	fmt.Println("Loading HA data...")
	entities, err := i.loadHomeAssistantEntities()
	if err != nil {
		return fmt.Errorf("failed to load ha data: %w", err)
	}
	fmt.Println("Parsing entities...")
	data := parse(entities)
	fmt.Println("Sending update to TRMNL...")
	if err := i.sendWebhook(data); err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	fmt.Println("Work done.")
	return nil
}
