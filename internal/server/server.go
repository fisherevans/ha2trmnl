package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"fisherevans.com/ha2trmnl/internal/config"
)

type PluginDataSource func() (map[string]any, error)

type Instance struct {
	config config.ServeConfig
	source PluginDataSource
}

func New(config config.ServeConfig, handler PluginDataSource) *Instance {
	if config.Address == "" {
		panic("missing address config")
	}
	if config.ApiToken == "" {
		panic("missing api token")
	}
	return &Instance{
		config: config,
		source: handler,
	}
}

const path = "/plugin_data"

func (i *Instance) MustStart() {
	slog.Info(fmt.Sprintf("listening on %s%s", i.config.Address, path))
	err := http.ListenAndServe(i.config.Address, i)
	if err != nil {
		panic(err)
	}
}

func (i *Instance) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) || strings.TrimPrefix(auth, prefix) != i.config.ApiToken {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		slog.Info("received unauthorized call")
		slog.Debug("unauthorized call details",
			"remote_addr", r.RemoteAddr,
			"method", r.Method,
			"url", r.URL.String(),
			"headers", r.Header,
		)
		return
	}

	if i.source == nil {
		slog.Error("no handler configured!")
		http.Error(w, "internal server error - server not configured correctly", http.StatusInternalServerError)
		return
	}

	resp, err := i.source()
	if err != nil {
		slog.Error("failed to handle request:", err)
		http.Error(w, "internal server error - failed to handle request", http.StatusInternalServerError)
		return
	}

	body, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		slog.Error("failed to marshal response:", err)
		http.Error(w, "internal server error - failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
	slog.Info("handled response")
	slog.Debug("response", "payload", resp)
}
