package configurer

import (
	"bytes"
	"encoding/json"
)

// A JSON configurer to use when the source is json encoded
type JSON struct {
	Source []byte
}

// Configure the item using the json-encoded source
func (h *JSON) Configure(x interface{}) error {
	return json.NewDecoder(bytes.NewBuffer(h.Source)).Decode(&x)
}
