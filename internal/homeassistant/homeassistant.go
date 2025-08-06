package homeassistant

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"fisherevans.com/ha2trmnl/internal/config"
	"github.com/gorilla/websocket"
)

type Instance struct {
	config config.HomeAssistantConfig
}

func New(config config.HomeAssistantConfig) *Instance {
	if config.HomeAssistantHost == "" {
		panic("missing ha host")
	}
	if config.HomeAssistantToken == "" {
		panic("missing ha token")
	}
	return &Instance{
		config: config,
	}
}

func (i *Instance) LoadHomeAssistantEntities() ([]Entity, error) {
	entities, err := fetchStates(i.config.HomeAssistantHost, i.config.HomeAssistantToken)
	if err != nil {
		return nil, fmt.Errorf("HA entity fetch: %w", err)
	}
	labelMap, err := fetchEntityLabelsWS(i.config.HomeAssistantHost, i.config.HomeAssistantToken)
	if err != nil {
		return nil, fmt.Errorf("HA label fetch: %w", err)
	}
	for i := range entities {
		if labels, ok := labelMap[entities[i].EntityID]; ok {
			entities[i].Labels = labels
		}
	}
	return entities, nil
}

func fetchStates(host, token string) ([]Entity, error) {
	url := fmt.Sprintf("http://%s/api/states", host)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HA error %d: %s", resp.StatusCode, body)
	}
	var arr []Entity
	return arr, json.Unmarshal(body, &arr)
}

func fetchEntityLabelsWS(host, token string) (map[string][]string, error) {
	url := fmt.Sprintf("ws://%s/api/websocket", host)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// Step 1: Authenticate
	authMsg := map[string]interface{}{
		"type":         "auth",
		"access_token": token,
	}
	if err := c.WriteJSON(authMsg); err != nil {
		return nil, err
	}

	// Wait for auth ok
	for {
		var msg map[string]interface{}
		if err := c.ReadJSON(&msg); err != nil {
			return nil, err
		}
		if msg["type"] == "auth_ok" {
			break
		}
		if msg["type"] == "auth_invalid" {
			return nil, fmt.Errorf("websocket auth failed: %v", msg)
		}
	}

	// Step 2: Request entity registry
	req := map[string]interface{}{
		"id":   1,
		"type": "config/entity_registry/list",
	}
	if err := c.WriteJSON(req); err != nil {
		return nil, err
	}

	// Step 3: Read response
	var response struct {
		Result []struct {
			EntityID string   `json:"entity_id"`
			Labels   []string `json:"labels"`
		} `json:"result"`
	}
	for {
		var msg map[string]interface{}
		if err := c.ReadJSON(&msg); err != nil {
			return nil, err
		}
		if msg["type"] == "result" && msg["id"].(float64) == 1 {
			raw, _ := json.Marshal(msg)
			_ = json.Unmarshal(raw, &response)
			break
		}
	}

	// Build label map
	labelMap := map[string][]string{}
	for _, e := range response.Result {
		if len(e.Labels) > 0 {
			labelMap[e.EntityID] = e.Labels
		}
	}
	return labelMap, nil
}
