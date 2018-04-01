package source

import "github.com/klokare/evo/config"

// Multi source combines multiple sources into one
type Multi []config.Source

// Value returns the first non-nil, if any, value from the composite sources.
func (m Multi) Value(ns []string, k string) interface{} {
	var x interface{}
	for _, s := range m {
		if x = s.Value(ns, k); x != nil {
			return x
		}
	}
	return nil
}
