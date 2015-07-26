package main

import (
	"html/template"
	"net/http"
	"runtime"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

var initTime = time.Now()

func init() {
	http.HandleFunc("/", handle)
	appengine.Main()
}

func main() {
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ctx := appengine.NewContext(r)
	log.Infof(ctx, "Serving the front page.")

	param := struct {
		RunningTime time.Duration
		NumCPU      int
	}{
		time.Since(initTime),
		runtime.NumCPU(),
	}

	tmpl.Execute(w, param)
}

var tmpl = template.Must(template.New("front").Parse(`
<html><body>

<p>
Hello, MVMs!
</p>

<p>
This instance has been running for <em>{{.RunningTime}}</em>.
</p>

<p>
Num of CPU is <em>{{.NumCPU}}</em>.
</p>

</body></html>
`))
