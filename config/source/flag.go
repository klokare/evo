package source

import (
	"flag"
	"strings"
)

// Flag provides a configuration source from the command-line flags
type Flag struct{}

// Value returns value of the flag which matches the key, if any. If not found, a second try is
// made using a lower-case version, substituting hypens for underscores
func (f Flag) Value(ns []string, k string) interface{} {

	// Rebuild the key including the namespace
	if len(ns) > 0 {
		k = strings.Join(ns, "-") + "-" + k
	}

	// Try as is
	if x := f.valueForKey(k); x != nil {
		return x
	}

	// Try using typical convention
	k = strings.ToLower(strings.Replace(k, "_", "-", -1))
	return f.valueForKey(k)
}

// Visit all flags that were actually set, looking for one with the same as the key
func (f Flag) valueForKey(k string) interface{} {
	var x interface{}
	flag.Visit(func(y *flag.Flag) {
		if y.Name == k {
			x = y.Value.(flag.Getter).Get()
		}
	})
	return x
}
