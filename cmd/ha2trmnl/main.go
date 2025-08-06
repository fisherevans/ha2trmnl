package main

import (
	"log/slog"
	"os"
	"time"
	_ "time/tzdata" // used to ensure timezone data is available if it's missing in the OS environment

	"fisherevans.com/ha2trmnl/internal/config"
	"fisherevans.com/ha2trmnl/internal/homeassistant"
	"fisherevans.com/ha2trmnl/internal/plugin_data"
	"fisherevans.com/ha2trmnl/internal/pusher"
	"fisherevans.com/ha2trmnl/internal/server"
)

func main() {
	if len(os.Args) <= 1 {
		panic("usage: ./ha2trmnl [serve|push|fetch] (./path/to/config.yaml)")
	}
	mode := os.Args[1]

	// load config
	configFile := "./config.yaml"
	if len(os.Args) >= 3 {
		configFile = os.Args[2]
	}
	c, err := config.FromFile(configFile)
	if err != nil {
		panic(err)
	}

	// setup logging
	logLevel := slog.LevelInfo
	if c.Debug {
		logLevel = slog.LevelDebug
	}
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(logHandler))

	// change local timezone based on config
	if c.TimeZone != "" {
		loc, err := time.LoadLocation(c.TimeZone)
		if err != nil {
			panic(err)
		}
		time.Local = loc
	}

	// setup dependencies
	slog.Info("loaded configuration from: " + configFile)
	ha := homeassistant.New(c.HomeAssistantConfig)
	source := plugin_data.New(ha.LoadHomeAssistantEntities)

	// run the app
	switch mode {
	case "serve":
		server.New(c.ServeConfig, source.Fetch).MustStart()
	case "push":
		pusher.New(c.PushConfig, source.Fetch).Run()
	case "fetch":
		fetch(source)
	default:
		panic("invalid mode!")
	}
}

func fetch(source *plugin_data.Instance) {
	data, err := source.Fetch()
	if err != nil {
		panic(err)
	}
	slog.Info("fetched plugin data", "data", data)
}
