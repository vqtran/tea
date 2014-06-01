package engines

/**
	Tea Engine Plugin for Mustache.
**/

import (
	"fmt"
	"github.com/hoisie/mustache"
	"io"
)

type MustacheEngine struct{}

// Mustache already provides a nice ParseFile function we just use that
// with default options.
func (m MustacheEngine) CompileFile(filepath string) (interface{}, error) {
	return mustache.ParseFile(filepath)
}

// Render the mustache template
func (m MustacheEngine) Render(buf io.Writer, tmpl interface{}, data interface{}) error {
	var r io.Writer
	rendered := tmpl.(*mustache.Template).Render(data, r)
	_, err := fmt.Fprintf(buf, rendered)
	return err
}

var Mustache = MustacheEngine{}
