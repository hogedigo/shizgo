package todo

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"net/http"
)

func init() {
	http.HandleFunc("/todo/done", delete)
}

func delete(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		http.Error(w, "login required.", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := datastore.DecodeKey(r.FormValue("key"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var todo Todo
	err = datastore.Get(c, key, &todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if u.ID != todo.UserId {
		http.Error(w, "forbidden access.", http.StatusForbidden)
		return
	}

	todo.Done = true

	_, err = datastore.Put(c, key, &todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/todo", http.StatusFound)
}
