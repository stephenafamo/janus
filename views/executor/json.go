package executor

import (
	"encoding/json"
	"io"
)

// JSON is a json executor that renders the data to json
type JSON struct{}

// Render renders data as JSON
func (j JSON) Render(wr io.Writer, name string, data interface{}) error {
	return json.NewEncoder(wr).Encode(data)
}

//Exists satisfies the executor interface
func (j JSON) Exists(name string) bool {
	// Every template exists
	return true
}

//Add satisfies the executor interface
func (j JSON) Add(name string, data io.Reader) error {
	// Does nothing.
	// Only here to satisfy the interface
	return nil
}
