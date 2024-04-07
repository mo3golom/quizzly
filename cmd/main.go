package main

import (
	"github.com/joncalhoun/form"
	"html/template"
	"net/http"
)

type User struct {
	Name     string
	UserName string
}

var inputTpl = `
<p><label {{with .ID}}for="{{.}}"{{end}}>
	{{.Label}}
</label>
<input {{with .ID}}id="{{.}}"{{end}} type="{{.Type}}" name="{{.Name}}" placeholder="{{.Placeholder}}" {{with .Value}}value="{{.}}"{{end}}></p>
{{with .Footer}}
  <p>{{.}}</p>
{{end}}
`

func main() {
	tpl := template.Must(template.New("").Parse(inputTpl))
	fb := form.Builder{
		InputTemplate: tpl,
	}

	pageTpl := template.Must(template.New("").Funcs(fb.FuncMap()).Parse(`
		<html>
		<body>
			<form>
				{{inputs_for .}}
			</form>
		</body>
		</html>`))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		pageTpl.Execute(w, User{})
	})
	http.ListenAndServe(":3000", nil)
}
