package util

import "encoding/json"

func ToJson(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
