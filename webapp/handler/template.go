package handler

import (
	"html/template"
	"net/http"
)

func init() {
	http.HandleFunc("/template", handleWithTemplate)
}

var src string = `
<html>
<body>
<a href="https://www.google.com/?q={{.Name}}">{{.Name}}</a><br>
{{$favorite := "orange"}} ʕ ◔ϖ◔ʔ .｡o(I love {{$favorite}}!)<br>
{{range .Fruits}}{{.}} {{if eq . $favorite}}yummy!{{end}}<br>{{end}}
<script>window.alert('{{.Name}}');</script>
</body>
</html>
`

func handleWithTemplate(w http.ResponseWriter, _ *http.Request) {

	t, err := template.New("test").Parse(src)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Name string
		Fruits []string
	}{
		"<b>'Hello' & Gopher!</b>",
		[]string{"apple", "banna", "orange"},
	}

	t.ExecuteTemplate(w, "test", data)
}
