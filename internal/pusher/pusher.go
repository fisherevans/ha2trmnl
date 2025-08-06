package pusher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"fisherevans.com/ha2trmnl/internal/config"
)

type PluginDataSource func() (map[string]any, error)

type Instance struct {
	config config.PushConfig
	source PluginDataSource
}

func New(config config.PushConfig, source PluginDataSource) *Instance {
	if config.Webhook == "" {
		panic("missing webhook url")
	}
	return &Instance{
		config: config,
		source: source,
	}
}

func (i *Instance) Run() {
	interval := time.Minute * 5
	if i.config.Interval != "" {
		var err error
		interval, err = time.ParseDuration(i.config.Interval)
		if err != nil {
			panic(err)
		}
	} else {
		slog.Warn("no interval specified, using 5 minute default")
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	slog.Info("starting periodic push", "interval", interval.String())

	for {
		i.runOnce()
		select {
		case <-ticker.C:
		}
	}
}

func (i *Instance) runOnce() {
	if err := i.sendData(); err != nil {
		slog.Error("failed to send plugin data:", err)
	} else {
		slog.Info("plugin data published.")
	}
}

func (i *Instance) sendData() error {
	pluginData, err := i.source()
	if err != nil {
		return fmt.Errorf("failed to load plugin data: %w", err)
	}
	body, err := json.Marshal(map[string]interface{}{
		"merge_variables": pluginData,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	slog.Info("sending plugin data: " + string(body))
	if i.config.DryRun {
		slog.Warn("DATA NOT SENT TO WEBHOOK DUE TO DRY RUN MODE")
		return nil
	}
	req, _ := http.NewRequest("POST", i.config.Webhook, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("error from webhook %d: %s", resp.StatusCode, respBody)
	}
	return nil
}
