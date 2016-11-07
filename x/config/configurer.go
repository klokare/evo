package config

import (
	"flag"
	"fmt"
	"reflect"
	"sync"
)

// A ConfigureError indicates that a problem occurred during configuration
type ConfigureError struct {
	Type     string
	Property string
	Message  string
	Success  bool
}

// Error returns the description of what went wrong during configuration
func (e ConfigureError) Error() string {
	s := "error configuring"
	if e.Type != "" {
		s += " " + e.Type
	}
	if e.Property != "" {
		s += "." + e.Property
	}
	s += ": " + e.Message
	return s
}

// Configured returns true if the configuration succeeded; otherwise, false.
func (e ConfigureError) Configured() bool { return e.Success }

// A Configurer can set the properties of a target struct using its internal settings
type Configurer struct {
	Settings interface{}
	once     sync.Once
	values   map[string]reflect.Value
}

// Configure the target interface using the stored settings. The target must be able to be set by
// reflection so pass a pointer to the helper. Configuration is explicit so the field in Settings
// must have the same "evo" struct tag as the field in the target.
func (h *Configurer) Configure(v interface{}) (err error) {

	// As reflection uses panics and not errors, wrap a recovery here to return the error instead
	defer func() {
		if r := recover(); r != nil {
			err = &ConfigureError{Message: fmt.Sprintf("%v", r)}
		}
	}()

	// Extract the properties from the settings and apply the command-line flags
	h.once.Do(func() { err = h.init() })
	if err != nil {
		return err
	}

	// There are no settings, return
	if h.Settings == nil || len(h.values) == 0 {
		return ConfigureError{Message: "no settings defined"}
	}

	// Transfer the settings to the target
	x := reflect.ValueOf(v)
	if x.Kind() != reflect.Struct {
		if x.Kind() != reflect.Ptr || reflect.ValueOf(x.Elem()).Kind() != reflect.Struct {
			return ConfigureError{Type: x.Type().Name(), Message: "target must be pointer to struct"}
		}
		x = x.Elem() // reference the struct
	}
	if !x.CanSet() {
		return ConfigureError{Type: x.Type().Name(), Message: "target cannot be changed"}
	}

	t := x.Type()
	for i := 0; i < t.NumField(); i++ {

		// This field is a struct or pointer to a struct
		y := x.Field(i)
		k := y.Kind()
		if k == reflect.Ptr || k == reflect.Interface {
			if y.Interface() != nil {
				if err = h.Configure(y); err != nil {
					return
				}
			}
		}

		// This field has an evo tag
		f := t.Field(i)
		e, ok := f.Tag.Lookup("evo")
		if !ok {
			continue // no evo tag for this field
		}

		// Update from setting
		z, ok := h.values[e]
		if !ok {
			continue // no matching setting
		}
		y.Set(z)
	}

	return nil
}

// TODO: Recursively check fields in settings if field is a struct or ptr to a struct
func (h *Configurer) init() (err error) {

	// As reflection uses panics and not errors, wrap a recovery here to return the error instead
	defer func() {
		if r := recover(); r != nil {
			err = &ConfigureError{Message: fmt.Sprintf("%v", r)}
		}
	}()

	// Ensure settings is the proper kind
	x := reflect.ValueOf(h.Settings)
	if x.Kind() == reflect.Ptr || x.Kind() == reflect.Interface {
		x = x.Elem()
	}
	if x.Kind() != reflect.Struct {
		return &ConfigureError{Type: x.Type().Name(), Message: "settings must be a struct or a pointer to a struct"}
	}

	// Extract the fields taged with "evo"
	t := x.Type()
	h.values = make(map[string]reflect.Value, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		e, ok := f.Tag.Lookup("evo")
		if !ok {
			continue
		}
		h.values[e] = x.Field(i)
	}

	// Override with command-line
	flag.Visit(func(f *flag.Flag) {
		if i, ok := Flags[f.Name]; ok {
			h.values[f.Name] = reflect.ValueOf(i).Elem()
		}
	})
	return nil
}
