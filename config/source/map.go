package source

// Map defines a source based on a standard map of interfaces with a string key.
type Map map[string]interface{}

// Value returns the value from the map source representd by the namespace and key; otherwise,
// returns nil.
func (m Map) Value(ns []string, k string) interface{} {

	// Iterate the namespaces
	var ok bool
	var x interface{}
	m1 := m
	for _, n := range ns {

		// Find the name in the current level
		if x, ok = m1[n]; !ok {
			return nil
		}

		// Value at name should be a map
		if m1, ok = x.(map[string]interface{}); !ok {
			return nil
		}
	}

	// Return the value from the map
	return m1[k]
}
