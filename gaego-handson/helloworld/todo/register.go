package todo

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"net/http"
)

func init() {
	http.HandleFunc("/todo/register", register)
}

type Todo struct {
	UserId  string
	Todo    string
	Notes   string
	DueDate string
	Done    bool
}

func register(w http.ResponseWriter, r *http.Request) {
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

	todo := Todo{
		u.ID,
		r.FormValue("Todo"),
		r.FormValue("Notes"),
		r.FormValue("DueDate"),
		false,
	}

	key := datastore.NewIncompleteKey(c, "Todo", nil)
	key, err := datastore.Put(c, key, &todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/todo", http.StatusFound)
}
