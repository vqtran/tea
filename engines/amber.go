package engines

/**
	Tea Engine Plugin for Amber.
**/

import (
	"github.com/eknkc/amber"
	"html/template"
	"io"
)

type AmberEngine struct{}

// Amber already provides a nice CompileFile function we just use that
// with default options.
func (a AmberEngine) CompileFile(filepath string) (interface{}, error) {
	return amber.CompileFile(filepath, amber.DefaultOptions)
}

func (a AmberEngine) Render(buf io.Writer, tmpl interface{}, data interface{}) error {
	return tmpl.(*template.Template).Execute(buf, data)
}

var Amber = AmberEngine{}
