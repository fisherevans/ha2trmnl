package runner

import (
	"fmt"
	"log"

	"fisherevans.com/ha2trmnl/internal/config"
	"fisherevans.com/ha2trmnl/internal/homeassistant"
	"fisherevans.com/ha2trmnl/internal/trmnl"
)

func Run(c config.Config) error {
	log.Println("Loading HA data...")
	entities, err := homeassistant.LoadHomeAssistantEntities(c)
	if err != nil {
		return fmt.Errorf("failed to load ha data: %w", err)
	}
	log.Println("Parsing entities...")
	data := parse(entities)
	log.Println("Sending update to TRMNL...")
	if err := trmnl.SendData(c, data); err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	log.Println("Work done.")
	return nil
}
