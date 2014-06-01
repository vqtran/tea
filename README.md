# Tea
*Lots of work in progress.*

Making Go templates easy.

Currently there are several different templating engines for Go, all with different usages specs, and if you want to load in multiple files or cache them you have to do it yourself. Tea does this for you.

Express.js has a great way of doing this:

```js
var app = express();

// Use EJS for templates
app.set('view engine', 'ejs');

// Look in app/views/ for the files
app.set('views', __dirname + '/app/views');

app.get('/', func(res, req) {
   // Render index with some data
   res.render('index', data);
});
```

Tea takes a similar approach:
```go
// Use Amber for templates
tea.SetEngine("amber")

// Look in templates/ for the files
tea.MustCompile("templates/", tea.Options{".amber", true})

http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
   // Get the compiled template in html/template form.
   tmpl, ok := tea.Get("index")
   if ok {
      tmpl.Execute(w, data)
   }
})
```

By doing this Tea provides several things:
   - Directory compilation (recursive or top-level)
   - Thread-safe template caching
   - Painless switch between templating engines
   - Compilation to the standard
   - Concise syntax
   - Extensible to support more engines.

## Documentation
```
package tea
    import "github.com/vqtran/tea"
```

### FUNCTIONS

#### func Clear()
    Clear the entire cache of templates.

#### func Compile(dirpath string, options Options) error
    Compile an entire directory with options.

#### func Delete(key string)
    Delete a specific key/value from the map.

#### func Get(key string) (*template.Template, bool)
    Thread-safe get the template. Second return variable states whether or
    not key is in the map.

#### func GetCache() map[string]*template.Template

#### func MustCompile(dirpath string, options Options)
    Same as Compile except will panic if there is an error

#### func SetEngine(en string) error
    Set what engine to use during compilation with a string. Also sets a
    default extension

### TYPES

```go
type Engine interface {
    // Has to be able to compile from a filepath to a html/template.
    CompileFile(filepath string) (*template.Template, error)
}
```
    Different templating engines must be support via a plugin in the engines
    subdirectory.


#### func GetEngine() *Engine
    Read-locked way to get the underlying tea-plugin for the engine. Makes
    it easy if user wants to compile a single file.

#### type Options struct {
    // File extension to match when compiling directories
    FileExt string
    // Whether search is recursive or only top-level
    Recursive bool
}
