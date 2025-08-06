package runner

import (
	"fmt"
	"log"
	"time"

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
	addGeneratedTimestamp(c, data)
	log.Println("Sending update to TRMNL...")
	if err := trmnl.SendData(c, data); err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	log.Println("Work done.")
	return nil
}

func addGeneratedTimestamp(c config.Config, m map[string]any) {
	t := time.Now()
	if c.TimeZone != "" {
		loc, err := time.LoadLocation(c.TimeZone)
		if err != nil {
			log.Println("failed to load timezone, defaulting to utc: ", err)
			t = t.In(time.UTC)
		} else {
			t = t.In(loc)
		}
	} else {
		t = t.In(time.Local)
	}
	m["generated"] = t.Format("15h04m")
}
