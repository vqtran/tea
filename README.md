# Tea
*Work in progress.*

Providing a simple way to use view templates in Go.

Currently there are several different templating engines for Go, all with different usages. Additionally, what if you wanted to compile and cache an entire directory of view templates.

Express.js has a great way of doing this:

```
var app = express();
app.set('view engine', 'ejs');
app.set('views', __dirname + '/app/views');
app.get('/', func(res, req) {
   res.render('index');
});
```

Tea takes a similar approach:
```go
tea.SetEngine("amber")
tea.Compile("templates/", options)
tmpl, err := tea.Get("index")
// tmpl is a *template.Template (html/template)
// ready to be executed.

```

By doing this Tea provides several things:
   - Directory compilation (recursive)
   - Template caching
   - Cache thread safety
   - Painless switching between templating engines
   - Concise syntax


## (Temporary) Documentation

type Engine interface {
   CompileFile(filepath)
}

tea.SetEngine("amber")
tea.GetEngine()
tea.Compile("directorypath", Options{recursive, extension, must})
tea.Get("index")
