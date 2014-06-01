package engines

/**
	Tea Engine Plugin for Amber.
**/

import (
	"github.com/eknkc/amber"
	"html/template"
)

type AmberEngine struct{}

// Amber already provides a nice CompileFile function we just use that
// with default options.
func (a AmberEngine) CompileFile(filepath string) (*template.Template, error) {
	return amber.CompileFile(filepath, amber.DefaultOptions)
}

var Amber = AmberEngine{}
