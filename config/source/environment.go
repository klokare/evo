package source

import (
	"os"
	"strings"
)

// Environment provides a configuration source from the environment variables
type Environment struct{}

// Value returns value of the environment variable which matches the key, if any. If not found, a
// second try is made using an upper-case version, substituting underscores for hypens.
func (e Environment) Value(ns []string, k string) interface{} {

	// Rebuild the key including the namespace
	if len(ns) > 0 {
		k = strings.Join(ns, "_") + "_" + k
	}

	// Try as is
	if x := e.valueForKey(k); x != nil {
		return x
	}

	// Try using typical convention
	k = strings.ToUpper(strings.Replace(k, "-", "_", -1))
	return e.valueForKey(k)
}

// Visit all flags that were actually set, looking for one with the same as the key
func (e Environment) valueForKey(k string) interface{} {
	if s, ok := os.LookupEnv(k); ok {
		return s
	}
	return nil
}
