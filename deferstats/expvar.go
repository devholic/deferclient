package deferstats

import (
	"expvar"
	"fmt"
)

// GetExpvar captures expvar using Do method
func (c *Client) GetExpvar() (string, error) {
	result := "{"
	first := true
	expvar.Do(func(kv expvar.KeyValue) {
		if !first {
			result += ",\n"
		}
		first = false
		result += fmt.Sprintf("%q: %s", kv.Key, kv.Value)
	})
	result += "}"

	return result, nil
}
