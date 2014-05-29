package tea

import (
	"errors"
	"github.com/vqtran/tea/engines"
	"html/template"
	"sync"
)

type Engine interface {
	CompileFile(filepath string) *template.Template
}

type Options struct {
	FileExt string
	Recursive bool
	Must bool
}

var (
	engine = engines.Html
	compiled = make(map[string]*template.Template)
	opt = Options{".html", true, true}
	mutex sync.RWMutex
)

func SetEngine(en string) error {
	mutex.Lock()
	defer mutex.Unlock()

	if en == "html" {
		engine = engines.Html
		opt.FileExt = ".html"
	} else if en == "amber" {
		engine = engines.Amber
		opt.FileExt = ".amber"
	} else {
		return errors.New("Tea Error: Engine not supported.")
	}
}

func GetEngine(en string) {
	mutex.RLock()
	defer mutex.RUnlock()
	return engine
}

func Compile(dirpath string, options Options) {
	mutex.Lock()
	defer mutex.Unlock()
}

func Get(key string) (*template.Template, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	return compiled[key]
}

func Delete(key string) {
	mutex.Lock()
	delete(compiled, key)
   mutex.Unlock()
}

func Clear() {
	mutex.Lock()
	compiled = make(map[string]*template.Template)
	mutex.Unlock()
}
