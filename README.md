# Tea
*Project discontinued, I would recommend [render](https://github.com/unrolled/render) for similar use cases.*

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

// Look in templates/ for the files ending in .amber, with recursive search on.
tea.MustCompile("templates/", tea.Options{".amber", true})

http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
   // Render index with some data
   err := tea.Render(w, "index", data)
   if err != nil {
      // Handle error
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

See example and documentation for more detailed usage.

## Supported Templating Engines
[html/template](http://golang.org/pkg/html/template/):
```go
tea.SetEngine("html")
```

[amber](https://github.com/eknkc/amber):
```go
tea.SetEngine("amber")
```

[mustache](https://github.com/hoisie/mustache)*
```go
tea.SetEngine("mustache")
```
*Mustache layouts currently do not work quite yet.

More to come!

## There's More
Plugins can be powerful. For example, Tea's implementation of the Go html/template plugin supports an "include" function/macro.

Before to use layouts you would do something like:

base.html
```html
<html>
   <head>
      <title>Hello World</title>
   </head>
   <body>
      {{ template "content" . }}
   </body>
</html>

```

content.html
```html
{{ define "content" }}
<h1>Hello World</h1>
<h2>Tea will knock your socks off.</h2>
{{ end }}
```

And you would parse it using.
```go
tmpl := template.Must(template.ParseFiles("content.html", "base.html"))
```

However, with Tea, writing `{{ include "file.html" }}` replaces that statement with the contents of file.html. So you can do:

content.html
```html
{{ define "content" }}
<h1>Hello World</h1>
<h2>Tea will knock your socks off.</h2>
{{ end }}
{{ include "base.html" }}
```

And Tea's compilation will then take care of it! Although this is less obvious that its a layout (versus normally stating at the top), this does make it clearer as to what the code is doing - and can be used for partials, anything where you want to include the contents of another file.

*NOTE: The macro right now is not recursive, being if the file you're including contains an include, it will fail. This prevents cycles, but if this is something that you think is very important then please post an issue and I'll look into it.*

## Contributing
There are tons of templating engines out there some - more than I could do quickly by myself. Write an engine plugin, test it, and send a pull request!

Remember, after you write a plugin edit SetEngine to support it. Current implementations can be found in `engines/`.

## Documentation
```go
package tea
    import "github.com/vqtran/tea"
```

### Functions

#### func Clear
```go
func Clear()
```
Clear the entire cache of templates.

#### func Compile
```go
func Compile(dirpath string, options Options) error
```
Compile an entire directory with options.

#### func Delete
```go
func Delete(key string)
```
Delete a specific key/value from the map.

#### func Get
```go
func Get(key string) (interface{}, bool)
```
Thread-safe get the template. Second return variable states whether or not key is in the map. Note that you will have to cast the template to use it. This can bee seen done in the implementation for the engines.

The lookup key is the path to the template relative to your specified templates directory without the extension. Therefore if you specified to look in "templates/" and in it you have "index.html" and "dir/example.html" you would access them with "index" and "dir/example" respectively.

#### func GetCache
```go
func GetCache() map[string]interface{}
```
Get the underlying hashmap for the cache.

#### func GetEngine
```go
func getEngine() *Engine
```
Read-locked way to get the underlying tea-plugin for the engine. Useful for if user needs to compile a single file.

#### func MustCompile
```go
func MustCompile(dirpath string, options Options)
```
Same as Compile except will panic if there is an error

#### func Render
```go
func Render(buf io.Writer, key string, data interface{}) error
```
Writes the rendered page with data to the io.Writer.

The lookup key is the path to the template relative to your specified templates directory without the extension. Therefore if you specified to look in "templates/" and in it you have "index.html" and "dir/example.html" you would access them with "index" and "dir/example" respectively.

#### func SetEngine
```go
func SetEngine(en string) error
```
Set what engine to use during compilation with a string. Also sets a default extension

### Types

```go
type Engine interface {
    // Has to be able to compile from a filepath
    CompileFile(filepath string) (*template.Template, error)
   // Function to render templates (does the casting)
   Render(buf io.Writer, tmpl interface{}, data interface{}) error
}
```
Different templating engines must be support via a plugin in the engines subdirectory.

```go
type Options struct {
    // File extension to match when compiling directories
    FileExt string
    // Whether search is recursive or only top-level
    Recursive bool
}
```
