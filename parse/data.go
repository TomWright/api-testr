package parse

import (
	"encoding/json"
	"fmt"
)

type data struct {
	d map[string]interface{}
}

func (d *data) UnmarshalJSON(data []byte) error {
	d.d = make(map[string]interface{})
	if err := json.Unmarshal(data, &d.d); err != nil {
		return fmt.Errorf("expected map[string]interface{}: %s", err)
	}

	return nil
}

func (d data) get(key string) (interface{}, bool) {
	val, ok := d.d[key]
	return val, ok
}

func (d data) string(key string) (string, bool) {
	val, ok := d.get(key)
	if !ok {
		return "", false
	}

	switch i := val.(type) {
	case string:
		return i, true
	case []byte:
		return string(i), true
	}

	return "", false
}

func (d data) int(key string) (int, bool) {
	val, ok := d.get(key)
	if !ok {
		return 0, false
	}

	switch i := val.(type) {
	case int:
		return i, true
	case int32:
		return int(i), true
	case int64:
		return int(i), true
	case float32:
		return int(i), true
	case float64:
		return int(i), true
	}

	return 0, false
}
