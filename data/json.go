package data

import (
	"encoding/json"
	"io"
)

// Serialize the given interface into a string based JSON format
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)

	return e.Encode(i)
}
