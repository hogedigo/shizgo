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
Hello {{.Name}}!
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
	}{
		"hogedigo",
	}

	t.ExecuteTemplate(w, "test", data)
}
