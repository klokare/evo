package config

import (
	"strconv"
	"strings"

	"github.com/klokare/evo"
)

// Source provides
type Source interface {
	Value(namespaces []string, key string) interface{}
}

// Configurer provides a query-based configuration helper
type Configurer struct {
	Source
}

// Int returns the int value for the key or 0.
func (c *Configurer) Int(key string) int {
	x := find(c.Source, key)
	if x == nil {
		return 0
	}
	if y, ok := x.(int); ok {
		return y
	}
	if y, ok := x.(string); ok {
		z, err := strconv.Atoi(y)
		if err == nil {
			return z
		}
	}
	if y, ok := x.(float64); ok {
		return int(y)
	}
	return 0
}

// Ints returns the slice of int values for the key or nil.
func (c *Configurer) Ints(key string) []int {
	x := find(c.Source, key)
	if x == nil {
		return nil
	}
	if y, ok := x.([]int); ok {
		return y
	}
	if y, ok := x.([]string); ok {
		var err error
		z := make([]int, len(y))
		for i := 0; i < len(y); i++ {
			z[i], err = strconv.Atoi(y[i])
			if err != nil {
				return nil
			}
		}
		return z
	}
	return nil
}

// Float64 returns the float64 value for the key or 0.0.
func (c *Configurer) Float64(key string) float64 {
	x := find(c.Source, key)
	if x == nil {
		return 0.0
	}
	if y, ok := x.(float64); ok {
		return y
	}
	if y, ok := x.(string); ok {
		z, err := strconv.ParseFloat(y, 64)
		if err == nil {
			return z
		}
	}
	return 0.0
}

// Float64s returns the slice of float64 values for the key or nil.
func (c *Configurer) Float64s(key string) []float64 {
	x := find(c.Source, key)
	if x == nil {
		return nil
	}
	if y, ok := x.([]float64); ok {
		return y
	}
	if y, ok := x.([]string); ok {
		var err error
		z := make([]float64, len(y))
		for i := 0; i < len(y); i++ {
			z[i], err = strconv.ParseFloat(y[i], 64)
			if err != nil {
				return nil
			}
		}
		return z
	}
	return nil
}

// Bool returns the boolean value for the key or false.
func (c *Configurer) Bool(key string) bool {
	x := find(c.Source, key)
	if x == nil {
		return false
	}
	if y, ok := x.(bool); ok {
		return y
	}
	if y, ok := x.(int); ok {
		return y != 0
	}
	if y, ok := x.(float64); ok {
		return int(y) != 0
	}
	if y, ok := x.(string); ok {
		return strings.ToLower(y) == "true"
	}
	return false
}

// Bools returns the slice of boolean values for the key or nil.
func (c *Configurer) Bools(key string) []bool {
	x := find(c.Source, key)
	if x == nil {
		return nil
	}
	if y, ok := x.([]bool); ok {
		return y
	}
	if y, ok := x.([]int); ok {
		z := make([]bool, len(y))
		for i := 0; i < len(y); i++ {
			z[i] = y[i] != 0
		}
		return z
	}
	if y, ok := x.([]string); ok {
		z := make([]bool, len(y))
		for i := 0; i < len(y); i++ {
			z[i] = strings.ToLower(y[i]) == "true"
		}
		return z
	}
	return nil
}

// String returns the string value for the key or "".
func (c *Configurer) String(key string) string {
	x := find(c.Source, key)
	if x == nil {
		return ""
	}
	if y, ok := x.(string); ok {
		return y
	}
	return ""
}

// Strings returns the slice of string values for the key or nil.
func (c *Configurer) Strings(key string) []string {
	x := find(c.Source, key)
	if x == nil {
		return nil
	}
	if y, ok := x.([]string); ok {
		return y
	}
	return nil
}

// Activation returns the activation for the key or 0.
func (c *Configurer) Activation(key string) evo.Activation {
	x := find(c.Source, key)
	if x == nil {
		return 0
	}
	if y, ok := x.(int); ok {
		return evo.Activation(y)
	}
	if y, ok := x.(float64); ok {
		return evo.Activation(int(y))
	}
	if y, ok := x.(string); ok {
		return evo.Activations[strings.ToLower(y)]
	}
	return 0
}

// Activations returns the slice of activations for the key or nil.
func (c *Configurer) Activations(key string) []evo.Activation {
	x := find(c.Source, key)
	if x == nil {
		return nil
	}
	if y, ok := x.([]int); ok {
		z := make([]evo.Activation, len(y))
		for i := 0; i < len(y); i++ {
			z[i] = evo.Activation(y[i])
		}
		return z
	}
	if y, ok := x.([]string); ok {
		z := make([]evo.Activation, len(y))
		for i := 0; i < len(y); i++ {
			z[i] = evo.Activations[y[i]]
		}
		return z
	}
	if y, ok := x.([]interface{}); ok {
		z := make([]evo.Activation, len(y))
		for i := 0; i < len(y); i++ {
			if a, ok := y[i].(int); ok {
				z[i] = evo.Activation(a)
			} else if b, ok := y[i].(string); ok {
				z[i] = evo.Activations[b]
			}
		}
		return z
	}
	return nil
}

// Comparison returns the comparison for the key or 0.
func (c *Configurer) Comparison(key string) evo.Comparison {
	x := find(c.Source, key)
	if x == nil {
		return 0
	}
	if y, ok := x.(int); ok {
		return evo.Comparison(y)
	}
	if y, ok := x.(float64); ok {
		return evo.Comparison(int(y))
	}
	if y, ok := x.(string); ok {
		return evo.Comparisons[strings.ToLower(y)]
	}
	return 0
}

// Comparisons returns the slice of comparisons for the key or nil.
func (c *Configurer) Comparisons(key string) []evo.Comparison {
	x := find(c.Source, key)
	if x == nil {
		return nil
	}
	if y, ok := x.([]int); ok {
		z := make([]evo.Comparison, len(y))
		for i := 0; i < len(y); i++ {
			z[i] = evo.Comparison(y[i])
		}
		return z
	}
	if y, ok := x.([]string); ok {
		z := make([]evo.Comparison, len(y))
		for i := 0; i < len(y); i++ {
			z[i] = evo.Comparisons[y[i]]
		}
		return z
	}
	if y, ok := x.([]interface{}); ok {
		z := make([]evo.Comparison, len(y))
		for i := 0; i < len(y); i++ {
			if a, ok := y[i].(int); ok {
				z[i] = evo.Comparison(a)
			} else if b, ok := y[i].(string); ok {
				z[i] = evo.Comparisons[b]
			}
		}
		return z
	}
	return nil
}

func split(key string) (ns []string, k string) {
	parts := strings.Split(key, "|")
	switch len(parts) {
	case 0, 1:
		k = key
		return
	default:
		ns = parts[:len(parts)-1]
		k = parts[len(parts)-1]
		return
	}
}

func find(src Source, key string) interface{} {
	// Find the value in the source
	var x interface{}
	ns, k := split(key)
	for i := len(ns); i > 0; i-- {
		x = src.Value(ns[:i], k)
		if x != nil {
			break // value found
		}
	}
	if x == nil {
		x = src.Value(nil, k) // try without a namespace
	}
	return x
}
