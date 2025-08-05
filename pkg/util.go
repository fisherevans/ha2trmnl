package pkg

import "encoding/json"

func toJson(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
