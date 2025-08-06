package plugin_data

import (
	"fmt"
	"log/slog"
	"time"

	"fisherevans.com/ha2trmnl/internal/homeassistant"
)

type EntityDataSource func() ([]homeassistant.Entity, error)

type Instance struct {
	source EntityDataSource
}

func New(source EntityDataSource) *Instance {
	return &Instance{
		source: source,
	}
}

func (i *Instance) Fetch() (map[string]any, error) {
	slog.Info("loading home assistant data...")
	entities, err := i.source()
	if err != nil {
		return nil, fmt.Errorf("failed to load ha data: %w", err)
	}
	data := parse(entities)
	data["generated"] = time.Now().In(time.Local).Format("15h04m")
	return data, nil
}
