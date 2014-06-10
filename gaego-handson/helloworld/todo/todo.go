package todo

import (
	"appengine"
	"appengine/user"
	"html/template"
	"net/http"
)

func init() {
	http.HandleFunc("/todo", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		loginUrl, _ := user.LoginURL(c, "/todo")
		http.Redirect(w, r, loginUrl, http.StatusFound)
		return
	}

	logoutUrl, _ := user.LogoutURL(c, "/")

	t, err := template.ParseFiles("todo/todo.tmpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "text/html; charset=utf-8")

	params := struct {
		LogoutUrl string
		User      *user.User
	}{
		logoutUrl,
		u,
	}

	err = t.Execute(w, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
