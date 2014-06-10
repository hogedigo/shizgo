package todo

import (
	"appengine"
	"appengine/user"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/todo", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		loginUrl, _ := user.LoginURL(c, "/todo")
		http.Redirect(w, r, loginUrl, http.StatusFound)
		return
	}
	logoutUrl, _ := user.LogoutURL(c, "/")

	html := `
<html><body>
Hello, %s ! - <a href="%s">sign out</a><br>
<hr>
This is TODO page under constuction!
</body></html>
`
	fmt.Fprintf(w, html, u.Email, logoutUrl)
}
