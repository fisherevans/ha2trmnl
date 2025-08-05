package pkg

type Entity struct {
	EntityID    string                 `json:"entity_id"`
	State       string                 `json:"state"`
	Attributes  map[string]interface{} `json:"attributes"`
	LastChanged string                 `json:"last_changed"`
	Labels      []string               `json:"labels"` // <-- added field
}

func (e Entity) stringAttribute(key, defaultValue string) string {
	val, ok := e.Attributes[key]
	if !ok {
		return defaultValue
	}
	stringVal, strOk := val.(string)
	if !strOk {
		return defaultValue
	}
	return stringVal
}

func (e Entity) HasAttributeValue(key string, desired any) bool {
	val, ok := e.Attributes[key]
	if !ok {
		return false
	}
	return val == desired
}

func (e Entity) FriendlyName() string {
	return e.stringAttribute("friendly_name", e.EntityID)
}

func (e Entity) DeviceClass() string {
	return e.stringAttribute("device_class", "")
}
