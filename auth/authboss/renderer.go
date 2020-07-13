package authboss

import (
	"bytes"
	"context"

	"github.com/stephenafamo/janus/views/executor"
	"github.com/volatiletech/authboss/v3"
)

// Renderer sastisfies the authboss.Renderer interface using the executor
type Renderer struct {
	Base      string
	Templates executor.Executor
}

// Load the listed templates
func (r Renderer) Load(names ...string) error {
	return nil
}

// Render a given template
func (r Renderer) Render(ctx context.Context, page string, data authboss.HTMLData) (output []byte, contentType string, err error) {
	var b bytes.Buffer

	err = r.Templates.Render(&b, page, data)
	if err != nil {
		return
	}

	return b.Bytes(), "text/html", err
}
