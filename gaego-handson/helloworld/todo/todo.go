package todo

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"encoding/json"
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

	q := datastore.NewQuery("Todo").Filter("UserId =", u.ID).Filter("Done =", false).Order("-DueDate")

	var todos []Todo
	keys, err := q.GetAll(c, &todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := struct {
		LogoutUrl string
		User      *user.User
		Todos     []Todo
		Keys      []*datastore.Key
	}{
		logoutUrl,
		u,
		todos,
		keys,
	}

	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-type", "application/json; charset=utf-8")
		b, err := json.MarshalIndent(todos, "", "\t")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	w.Header().Set("Content-type", "text/html; charset=utf-8")

	err = t.Execute(w, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
