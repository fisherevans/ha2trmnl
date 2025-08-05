package trmnl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"fisherevans.com/ha2trmnl/internal/config"
)

func SendData(c config.Config, data map[string]interface{}) error {
	payload := map[string]interface{}{
		"merge_variables": data,
	}
	body, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	log.Println("Payload sent to TRMNL:\n" + string(body))
	if c.DryRun {
		log.Println("!!! DATA NOT SENT TO WEBHOOK DUE TO DRY RUN MODE !!!")
		return nil
	}
	req, _ := http.NewRequest("POST", c.TrmnlWebhook, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("TRMNL error %d: %s", resp.StatusCode, respBody)
	}
	return nil
}
