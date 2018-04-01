package config

import (
	"strings"
	"testing"

	"github.com/klokare/evo"
)

var (
	src mockSource = map[string]interface{}{
		// string
		"string-a":             "foo",
		"name1|string-b":       "goo",
		"name2|name3|string-c": "hoo",
		"string-bad":           123,

		// int
		"int-a":             1,
		"name1|int-b":       "2", // value may come in as string
		"name2|name3|int-c": 3,
		"int-bad":           "bad",

		// float64
		"float-a":             1.1,
		"name1|float-b":       "2.2", // value may come in a string
		"name2|name3|float-c": 3.3,
		"float-bad":           "bad",

		// bool
		"bool-a":              true,
		"name1|bool-b1":       0,       // can use 0 = false; nonzero = true, too
		"name1|bool-b2":       1,       // can use 0 = false; nonzero = true, too
		"name2|name3|bool-c1": "false", // value may come in as string
		"name2|name3|bool-c2": "true",  // value may come in as string
		"bool-bad":            "foo",

		// activation
		"activation-a":             1,
		"name1|activation-b":       "sigmoid", // value may come in a string
		"name2|name3|activation-c": 3,
		"activation-bad":           "foo",

		// compare -- only available as string
		"compare-a":             1,
		"name1|compare-b":       "age", // value may come in a string
		"name2|name3|compare-c": "novelty",
		"compare-bad":           "foo",

		// ints
		"ints-a":             []int{1, 2},
		"name1|ints-b":       []string{"2", "3"}, // value may come in as string
		"name2|name3|ints-c": []int{3, 4},
		"ints-bad":           []string{"5", "bad"},

		// strings
		"strings-a":             []string{"foo", "fee"},
		"name1|strings-b":       []string{"goo", "gee"},
		"name2|name3|strings-c": []string{"hoo", "hee"},
		"strings-bad":           []int{123, 234},

		// float64s
		"floats-a":             []float64{1.1, 2.2},
		"name1|floats-b":       []string{"1.1", "2.2"}, // value may come in a string
		"name2|name3|floats-c": []float64{3.3, 4.4},
		"floats-bad":           []string{"1.23", "bad"},

		// bools
		"bools-a":             []bool{false, true},
		"name1|bools-b":       []int{0, 1},               // can use 0 = false; nonzero = true, too
		"name2|name3|bools-c": []string{"false", "true"}, // value may come in as string
		"bools-bad":           []float64{1.1, 2.2},

		// activations
		"activations-a":             []int{1, 2},
		"name1|activations-b":       []string{"sigmoid", "tanh"}, // value may come in a string
		"name2|name3|activations-c": []int{3, 4},
		"activations-bad":           []float64{1.1, 2.2},

		// compare -- only available as string
		"compares-a":             []int{1, 2},
		"name1|compares-b":       []string{"age", "novelty"}, // value may come in a string
		"name2|name3|compares-c": []string{"novelty", "complexity"},
		"compares-bad":           []float64{1.1, 2.2},
	}
)

type mockSource map[string]interface{}

func (m mockSource) Value(namespaces []string, key string) interface{} {
	x := make([]string, 0, len(namespaces)+1)
	x = append(x, namespaces...)
	x = append(x, key)
	return m[strings.Join(x, "|")]
}

func TestInt(t *testing.T) {
	var cases = []struct {
		Desc     string
		Key      string
		Expected int
	}{
		{
			Desc:     "unknown key",
			Key:      "int-unknown",
			Expected: 0,
		},
		{
			Desc:     "not an int",
			Key:      "int-bad",
			Expected: 0,
		},
		{
			Desc:     "no namespace",
			Key:      "int-a",
			Expected: 1,
		},
		{
			Desc:     "one namespace",
			Key:      "name1|int-b",
			Expected: 2,
		},
		{
			Desc:     "one namespace, key found at higher level",
			Key:      "name1|int-a",
			Expected: 1,
		},
		{
			Desc:     "two namespaces",
			Key:      "name2|name3|int-c",
			Expected: 3,
		},
		{
			Desc:     "two namespaces, key found at mid level",
			Key:      "name1|name2|int-b",
			Expected: 2,
		},
		{
			Desc:     "two namespaces, key found at highest level",
			Key:      "name1|name2|int-a",
			Expected: 1,
		},
	}

	// Create the configurer
	cfg := &Configurer{Source: src}

	// Iterate the cases
	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			x := cfg.Int(c.Key)
			if c.Expected != x {
				t.Errorf("incorrect response for %s: expected %d, actual %d", c.Key, c.Expected, x)
			}
		})
	}
}

func TestInts(t *testing.T) {
	var cases = []struct {
		Desc     string
		Key      string
		Expected []int
	}{
		{
			Desc:     "unknown key",
			Key:      "ints-unknown",
			Expected: nil,
		},
		{
			Desc:     "not an int",
			Key:      "ints-bad",
			Expected: nil,
		},
		{
			Desc:     "no namespace",
			Key:      "ints-a",
			Expected: []int{1, 2},
		},
		{
			Desc:     "one namespace",
			Key:      "name1|ints-b",
			Expected: []int{2, 3},
		},
		{
			Desc:     "one namespace, key found at higher level",
			Key:      "name1|ints-a",
			Expected: []int{1, 2},
		},
		{
			Desc:     "two namespaces",
			Key:      "name2|name3|ints-c",
			Expected: []int{3, 4},
		},
		{
			Desc:     "two namespaces, key found at mid level",
			Key:      "name1|name2|ints-b",
			Expected: []int{2, 3},
		},
		{
			Desc:     "two namespaces, key found at highest level",
			Key:      "name1|name2|ints-a",
			Expected: []int{1, 2},
		},
	}

	// Create the configurer
	cfg := &Configurer{Source: src}

	// Iterate the cases
	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			x := cfg.Ints(c.Key)
			if len(c.Expected) != len(x) {
				t.Errorf("incorrect value length for %s: expected %d, actual %d", c.Key, len(c.Expected), len(x))
			} else {
				for i, e := range c.Expected {
					if e != x[i] {
						t.Errorf("incorrect value at %d for %s: expected %d, actual %d", i, c.Key, e, x[i])
					}
				}
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	var cases = []struct {
		Desc     string
		Key      string
		Expected float64
	}{
		{
			Desc:     "unknown key",
			Key:      "float-unknown",
			Expected: 0.0,
		},
		{
			Desc:     "not an float",
			Key:      "float-bad",
			Expected: 0.0,
		},
		{
			Desc:     "no namespace",
			Key:      "float-a",
			Expected: 1.1,
		},
		{
			Desc:     "one namespace",
			Key:      "name1|float-b",
			Expected: 2.2,
		},
		{
			Desc:     "one namespace, key found at higher level",
			Key:      "name1|float-a",
			Expected: 1.1,
		},
		{
			Desc:     "two namespaces",
			Key:      "name2|name3|float-c",
			Expected: 3.3,
		},
		{
			Desc:     "two namespaces, key found at mid level",
			Key:      "name1|name2|float-b",
			Expected: 2.2,
		},
		{
			Desc:     "two namespaces, key found at highest level",
			Key:      "name1|name2|float-a",
			Expected: 1.1,
		},
	}

	// Create the configurer
	cfg := &Configurer{Source: src}

	// Iterate the cases
	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			x := cfg.Float64(c.Key)
			if c.Expected != x {
				t.Errorf("incorrect response for %s: expected %f, actual %f", c.Key, c.Expected, x)
			}
		})
	}
}

func TestFloat64s(t *testing.T) {
	var cases = []struct {
		Desc     string
		Key      string
		Expected []float64
	}{
		{
			Desc:     "unknown key",
			Key:      "floats-unknown",
			Expected: nil,
		},
		{
			Desc:     "not an floats",
			Key:      "floats-bad",
			Expected: nil,
		},
		{
			Desc:     "no namespace",
			Key:      "floats-a",
			Expected: []float64{1.1, 2.2},
		},
		{
			Desc:     "one namespace",
			Key:      "name1|floats-b",
			Expected: []float64{1.1, 2.2},
		},
		{
			Desc:     "one namespace, key found at higher level",
			Key:      "name1|floats-a",
			Expected: []float64{1.1, 2.2},
		},
		{
			Desc:     "two namespaces",
			Key:      "name2|name3|floats-c",
			Expected: []float64{3.3, 4.4},
		},
		{
			Desc:     "two namespaces, key found at mid level",
			Key:      "name1|name2|floats-b",
			Expected: []float64{1.1, 2.2},
		},
		{
			Desc:     "two namespaces, key found at highest level",
			Key:      "name1|name2|floats-a",
			Expected: []float64{1.1, 2.2},
		},
	}

	// Create the configurer
	cfg := &Configurer{Source: src}

	// Iterate the cases
	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			x := cfg.Float64s(c.Key)
			if len(c.Expected) != len(x) {
				t.Errorf("incorrect value length for %s: expected %d, actual %d", c.Key, len(c.Expected), len(x))
			} else {
				for i, e := range c.Expected {
					if e != x[i] {
						t.Errorf("incorrect value at %d for %s: expected %f, actual %f", i, c.Key, e, x[i])
					}
				}
			}
		})
	}
}
func TestBool(t *testing.T) {
	var cases = []struct {
		Desc     string
		Key      string
		Expected bool
	}{
		{
			Desc:     "unknown key",
			Key:      "bool-unknown",
			Expected: false,
		},
		{
			Desc:     "not a bool",
			Key:      "bool-bad",
			Expected: false,
		},
		{
			Desc:     "no namespace",
			Key:      "bool-a",
			Expected: true,
		},
		{
			Desc:     "one namespace",
			Key:      "name1|bool-b1",
			Expected: false,
		},
		{
			Desc:     "one namespace",
			Key:      "name1|bool-b2",
			Expected: true,
		},
		{
			Desc:     "one namespace, key found at higher level",
			Key:      "name1|bool-a",
			Expected: true,
		},
		{
			Desc:     "two namespaces",
			Key:      "name2|name3|bool-c1",
			Expected: false,
		},
		{
			Desc:     "two namespaces",
			Key:      "name2|name3|bool-c2",
			Expected: true,
		},
		{
			Desc:     "two namespaces, key found at mid level",
			Key:      "name1|name2|bool-b2",
			Expected: true,
		},
		{
			Desc:     "two namespaces, key found at highest level",
			Key:      "name1|name2|bool-a",
			Expected: true,
		},
	}

	// Create the configurer
	cfg := &Configurer{Source: src}

	// Iterate the cases
	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			x := cfg.Bool(c.Key)
			if c.Expected != x {
				t.Errorf("incorrect response for %s: expected %t, actual %t", c.Key, c.Expected, x)
			}
		})
	}
}

func TestBools(t *testing.T) {
	var cases = []struct {
		Desc     string
		Key      string
		Expected []bool
	}{
		{
			Desc:     "unknown key",
			Key:      "bools-unknown",
			Expected: nil,
		},
		{
			Desc:     "not an bool",
			Key:      "bools-bad",
			Expected: nil,
		},
		{
			Desc:     "no namespace",
			Key:      "bools-a",
			Expected: []bool{false, true},
		},
		{
			Desc:     "one namespace",
			Key:      "name1|bools-b",
			Expected: []bool{false, true},
		},
		{
			Desc:     "one namespace, key found at higher level",
			Key:      "name1|bools-a",
			Expected: []bool{false, true},
		},
		{
			Desc:     "two namespaces",
			Key:      "name2|name3|bools-c",
			Expected: []bool{false, true},
		},
		{
			Desc:     "two namespaces, key found at mid level",
			Key:      "name1|name2|bools-b",
			Expected: []bool{false, true},
		},
		{
			Desc:     "two namespaces, key found at highest level",
			Key:      "name1|name2|bools-a",
			Expected: []bool{false, true},
		},
	}

	// Create the configurer
	cfg := &Configurer{Source: src}

	// Iterate the cases
	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			x := cfg.Bools(c.Key)
			if len(c.Expected) != len(x) {
				t.Errorf("incorrect value length for %s: expected %d, actual %d", c.Key, len(c.Expected), len(x))
			} else {
				for i, e := range c.Expected {
					if e != x[i] {
						t.Errorf("incorrect value at %d for %s: expected %t, actual %t", i, c.Key, e, x[i])
					}
				}
			}
		})
	}
}

func TestString(t *testing.T) {
	var cases = []struct {
		Desc     string
		Key      string
		Expected string
	}{
		{
			Desc:     "unknown key",
			Key:      "string-unknown",
			Expected: "",
		},
		{
			Desc:     "not a string",
			Key:      "string-bad",
			Expected: "",
		},
		{
			Desc:     "no namespace",
			Key:      "string-a",
			Expected: "foo",
		},
		{
			Desc:     "one namespace",
			Key:      "name1|string-b",
			Expected: "goo",
		},
		{
			Desc:     "one namespace, key found at higher level",
			Key:      "name1|string-a",
			Expected: "foo",
		},
		{
			Desc:     "two namespaces",
			Key:      "name2|name3|string-c",
			Expected: "hoo",
		},
		{
			Desc:     "two namespaces, key found at mid level",
			Key:      "name1|name2|string-b",
			Expected: "goo",
		},
		{
			Desc:     "two namespaces, key found at highest level",
			Key:      "name1|name2|string-a",
			Expected: "foo",
		},
	}

	// Create the configurer
	cfg := &Configurer{Source: src}

	// Iterate the cases
	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			x := cfg.String(c.Key)
			if c.Expected != x {
				t.Errorf("incorrect response for %s: expected %s, actual %s", c.Key, c.Expected, x)
			}
		})
	}
}

func TestStrings(t *testing.T) {
	var cases = []struct {
		Desc     string
		Key      string
		Expected []string
	}{
		{
			Desc:     "unknown key",
			Key:      "strings-unknown",
			Expected: nil,
		},
		{
			Desc:     "not a string",
			Key:      "strings-bad",
			Expected: nil,
		},
		{
			Desc:     "no namespace",
			Key:      "strings-a",
			Expected: []string{"foo", "fee"},
		},
		{
			Desc:     "one namespace",
			Key:      "name1|strings-b",
			Expected: []string{"goo", "gee"},
		},
		{
			Desc:     "one namespace, key found at higher level",
			Key:      "name1|strings-a",
			Expected: []string{"foo", "fee"},
		},
		{
			Desc:     "two namespaces",
			Key:      "name2|name3|strings-c",
			Expected: []string{"hoo", "hee"},
		},
		{
			Desc:     "two namespaces, key found at mid level",
			Key:      "name1|name2|strings-b",
			Expected: []string{"goo", "gee"},
		},
		{
			Desc:     "two namespaces, key found at highest level",
			Key:      "name1|name2|strings-a",
			Expected: []string{"foo", "fee"},
		},
	}

	// Create the configurer
	cfg := &Configurer{Source: src}

	// Iterate the cases
	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			x := cfg.Strings(c.Key)
			if len(c.Expected) != len(x) {
				t.Errorf("incorrect value length for %s: expected %d, actual %d", c.Key, len(c.Expected), len(x))
			} else {
				for i, e := range c.Expected {
					if e != x[i] {
						t.Errorf("incorrect value at %d for %s: expected %s, actual %s", i, c.Key, e, x[i])
					}
				}
			}
		})
	}
}

func TestActivation(t *testing.T) {
	var cases = []struct {
		Desc     string
		Key      string
		Expected evo.Activation
	}{
		{
			Desc:     "unknown key",
			Key:      "activation-unknown",
			Expected: 0,
		},
		{
			Desc:     "not an activation",
			Key:      "activation-bad",
			Expected: 0,
		},
		{
			Desc:     "no namespace",
			Key:      "activation-a",
			Expected: evo.Direct,
		},
		{
			Desc:     "one namespace",
			Key:      "name1|activation-b",
			Expected: evo.Sigmoid,
		},
		{
			Desc:     "one namespace, key found at higher level",
			Key:      "name1|activation-a",
			Expected: evo.Direct,
		},
		{
			Desc:     "two namespaces",
			Key:      "name2|name3|activation-c",
			Expected: evo.SteepenedSigmoid,
		},
		{
			Desc:     "two namespaces, key found at mid level",
			Key:      "name1|name2|activation-b",
			Expected: evo.Sigmoid,
		},
		{
			Desc:     "two namespaces, key found at highest level",
			Key:      "name1|name2|activation-a",
			Expected: evo.Direct,
		},
	}

	// Create the configurer
	cfg := &Configurer{Source: src}

	// Iterate the cases
	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			x := cfg.Activation(c.Key)
			if c.Expected != x {
				t.Errorf("incorrect response for %s: expected %s, actual %s", c.Key, c.Expected.String(), x.String())
			}
		})
	}
}

func TestActivations(t *testing.T) {
	var cases = []struct {
		Desc     string
		Key      string
		Expected []evo.Activation
	}{
		{
			Desc:     "unknown key",
			Key:      "activations-unknown",
			Expected: nil,
		},
		{
			Desc:     "not a string",
			Key:      "activations-bad",
			Expected: nil,
		},
		{
			Desc:     "no namespace",
			Key:      "activations-a",
			Expected: []evo.Activation{evo.Direct, evo.Sigmoid},
		},
		{
			Desc:     "one namespace",
			Key:      "name1|activations-b",
			Expected: []evo.Activation{evo.Sigmoid, evo.Tanh},
		},
		{
			Desc:     "one namespace, key found at higher level",
			Key:      "name1|activations-a",
			Expected: []evo.Activation{evo.Direct, evo.Sigmoid},
		},
		{
			Desc:     "two namespaces",
			Key:      "name2|name3|activations-c",
			Expected: []evo.Activation{evo.SteepenedSigmoid, evo.Tanh},
		},
		{
			Desc:     "two namespaces, key found at mid level",
			Key:      "name1|name2|activations-b",
			Expected: []evo.Activation{evo.Sigmoid, evo.Tanh},
		},
		{
			Desc:     "two namespaces, key found at highest level",
			Key:      "name1|name2|activations-a",
			Expected: []evo.Activation{evo.Direct, evo.Sigmoid},
		},
	}

	// Create the configurer
	cfg := &Configurer{Source: src}

	// Iterate the cases
	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			t.Run(c.Desc, func(t *testing.T) {
				x := cfg.Activations(c.Key)
				if len(c.Expected) != len(x) {
					t.Errorf("incorrect value length for %s: expected %d, actual %d", c.Key, len(c.Expected), len(x))
				} else {
					for i, e := range c.Expected {
						if e != x[i] {
							t.Errorf("incorrect value at %d for %s: expected %s, actual %s", i, c.Key, e, x[i])
						}
					}
				}
			})
		})
	}
}

func TestComparison(t *testing.T) {
	var cases = []struct {
		Desc     string
		Key      string
		Expected evo.Comparison
	}{
		{
			Desc:     "unknown key",
			Key:      "compare-unknown",
			Expected: 0,
		},
		{
			Desc:     "not a comparison",
			Key:      "compare-bad",
			Expected: 0,
		},
		{
			Desc:     "no namespace",
			Key:      "compare-a",
			Expected: evo.ByFitness,
		},
		{
			Desc:     "one namespace",
			Key:      "name1|compare-b",
			Expected: evo.ByAge,
		},
		{
			Desc:     "one namespace, key found at higher level",
			Key:      "name1|compare-a",
			Expected: evo.ByFitness,
		},
		{
			Desc:     "two namespaces",
			Key:      "name2|name3|compare-c",
			Expected: evo.ByNovelty,
		},
		{
			Desc:     "two namespaces, key found at mid level",
			Key:      "name1|name2|compare-b",
			Expected: evo.ByAge,
		},
		{
			Desc:     "two namespaces, key found at highest level",
			Key:      "name1|name2|compare-a",
			Expected: evo.ByFitness,
		},
	}

	// Create the configurer
	cfg := &Configurer{Source: src}

	// Iterate the cases
	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			x := cfg.Comparison(c.Key)
			if c.Expected != x {
				t.Errorf("incorrect response for %s: expected %s, actual %s", c.Key, c.Expected.String(), x.String())
			}
		})
	}
}

func TestComparisons(t *testing.T) {
	var cases = []struct {
		Desc     string
		Key      string
		Expected []evo.Comparison
	}{
		{
			Desc:     "unknown key",
			Key:      "compares-unknown",
			Expected: nil,
		},
		{
			Desc:     "not a string",
			Key:      "compares-bad",
			Expected: nil,
		},
		{
			Desc:     "no namespace",
			Key:      "compares-a",
			Expected: []evo.Comparison{evo.ByFitness, evo.ByNovelty},
		},
		{
			Desc:     "one namespace",
			Key:      "name1|compares-b",
			Expected: []evo.Comparison{evo.ByAge, evo.ByNovelty},
		},
		{
			Desc:     "one namespace, key found at higher level",
			Key:      "name1|compares-a",
			Expected: []evo.Comparison{evo.ByFitness, evo.ByNovelty},
		},
		{
			Desc:     "two namespaces",
			Key:      "name2|name3|compares-c",
			Expected: []evo.Comparison{evo.ByNovelty, evo.ByComplexity},
		},
		{
			Desc:     "two namespaces, key found at mid level",
			Key:      "name1|name2|compares-b",
			Expected: []evo.Comparison{evo.ByAge, evo.ByNovelty},
		},
		{
			Desc:     "two namespaces, key found at highest level",
			Key:      "name1|name2|compares-a",
			Expected: []evo.Comparison{evo.ByFitness, evo.ByNovelty},
		},
	}

	// Create the configurer
	cfg := &Configurer{Source: src}

	// Iterate the cases
	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			t.Run(c.Desc, func(t *testing.T) {
				x := cfg.Comparisons(c.Key)
				if len(c.Expected) != len(x) {
					t.Errorf("incorrect value length for %s: expected %d, actual %d", c.Key, len(c.Expected), len(x))
				} else {
					for i, e := range c.Expected {
						if e != x[i] {
							t.Errorf("incorrect value at %d for %s: expected %s, actual %s", i, c.Key, e, x[i])
						}
					}
				}
			})
		})
	}
}
