package config

import (
	"flag"
	"log"
	"reflect"
	"strings"
)

// Settings is the collection of properties from all the helpers in the EVO library. It provides
// a convenient place, especially when used with config.Configurer, to manage these properies
// in a new experiment.
type Settings struct {

	// Trial runner
	Trials int `evo:"trials" json:"trials" toml:"trials"`

	// Experiment
	Description string `evo:"description" json:"description" toml:"description"`
	Iterations  int    `evo:"iterations" json:"iterations" toml:"iterations"`

	// Distance comparer
	ConnsCoefficient  float64 `evo:"conns-coefficient" json:"conns-coefficient" toml:"conns-coefficient"`
	NodesCoefficient  float64 `evo:"nodes-coefficient" json:"nodes-coefficient" toml:"nodes-coefficient"`
	WeightCoefficient float64 `evo:"weight-coefficient" json:"weight-coefficient" toml:"weight-coefficient"`

	// Multiple crosser
	EnableProbability float64 `evo:"enable-probability" json:"enable-probability" toml:"enable-probability"`

	// Complexify mutator
	AddNodeProbability float64 `evo:"add-node-probability" json:"add-node-probability" toml:"add-node-probability"`
	AddConnProbability float64 `evo:"add-conn-probability" json:"add-conn-probability" toml:"add-conn-probability"`
	AllowRecurrent     bool    `evo:"allow-recurrent" json:"allow-recurrent" toml:"allow-recurrent"`

	// Weight mutator
	MutateWeightProbability  float64 `evo:"mutate-weight-probability" json:"mutate-weight-probability" toml:"mutate-weight-probability"`
	ReplaceWeightProbability float64 `evo:"replace-weight-probability" json:"replace-weight-probability" toml:"replace-weight-probability"`

	// Generational searcher
	SurvivalRate                float64 `evo:"survival-rate" json:"survival-rate" toml:"survival-rate"`
	MaxStagnation               int     `evo:"max-stagnation" json:"max-stagnation" toml:"max-stagnation"`
	MutateOnlyProbability       float64 `evo:"mutate-only-probability" json:"mutate-only-probability" toml:"mutate-only-probability"`
	InterspeciesMateProbability float64 `evo:"interspecies-mate-probability" json:"interspecies-mate-probability" toml:"interspecies-mate-probability"`

	// Static speciator
	CompatibilityThreshold float64 `evo:"compatibility-threshold" json:"compatibility-threshold" toml:"compatibility-threshold"`

	// Dynamic speciator
	TargetSpecies         int     `evo:"target-species" json:"target-species" toml:"target-species"`
	CompatibilityModifier float64 `evo:"compatibility-modifier" json:"compatibility-modifier" toml:"compatibility-modifier"`

	// NEAT seeder
	PopulationSize int `evo:"population-size" json:"population-size" toml:"population-size"`
	NumInputs      int `evo:"num-inputs" json:"num-inputs" toml:"num-inputs"`
	NumOutputs     int `evo:"num-outputs" json:"num-outputs" toml:"num-outputs"`

	// Web watcher
	WebURL          string `evo:"web-url" json:"web-url" toml:"web-url"`
	WebStudyID      int    `evo:"web-study-id" json:"web-study-id" toml:"web-study-id"`
	WebExperimentID int    `evo:"web-experiment-id" json:"web-experiment-id" toml:"web-experiment-id"`
}

// Flags are the command-line flags derived from the settings object
var Flags map[string]interface{}

func init() {
	Flags = make(map[string]interface{})
	SetFlags(Settings{})
}

// SetFlags creates the command-line flags based on the evo-tagged fields in the struct passed as an argument
func SetFlags(v interface{}) {

	// Ensure we are working with a struct
	x := reflect.ValueOf(v)
	if x.Kind() != reflect.Struct {
		x = x.Elem()
		if x.Kind() != reflect.Struct {
			log.Println("INFO", "config.SetFlags", reflect.ValueOf(v).Kind(), "cannot be treated as struct")
			return
		}
	}

	// Iterate the fields
	t := x.Type()
	for i := 0; i < t.NumField(); i++ {
		e, ok := t.Field(i).Tag.Lookup("evo")
		if !ok {
			continue
		}

		if _, ok = Flags[e]; !ok {
			f := t.Field(i)
			switch f.Type.Kind() {
			case reflect.String:
				Flags[e] = flag.String(e, "", strings.Replace(e, "-", " ", -1))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				Flags[e] = flag.Int(e, 0, strings.Replace(e, "-", "", -1))
			case reflect.Bool:
				Flags[e] = flag.Bool(e, false, strings.Replace(e, "-", "", -1))
			case reflect.Float32, reflect.Float64:
				Flags[e] = flag.Float64(e, 0.0, strings.Replace(e, "-", "", -1))
			}
		}
	}
}
