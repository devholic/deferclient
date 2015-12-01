package expvars

import (
	"expvar"
	"fmt"
)

// GetByDo takes expvar map using Do method from expvar package
func GetByDo() (string, error) {
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
