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

## (Temporary) Documentation
```
type Engine interface {
   CompileFile(filepath)
}

tea.SetEngine("amber")
tea.GetEngine()
tea.Compile("directorypath", Options{recursive, extension, must})
tea.Get("index")
```
