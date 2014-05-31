package engines

import (
	"github.com/eknkc/amber"
	"html/template"
)

type AmberEngine struct{}

func (a AmberEngine) CompileFile(filepath string) (*template.Template, error) {
	return amber.CompileFile(filepath, amber.DefaultOptions)
}

var Amber = AmberEngine{}
