package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"code.google.com/p/go.net/html"
)

func init() {
	http.HandleFunc("/pcrawl", handlePcrawl)
}

type syncUrlset struct {
	set map[string]bool
	sync.Mutex
}

func newSyncUrlset() *syncUrlset {
	var urlset syncUrlset
	urlset.set = make(map[string]bool)
	return &urlset
}

func handlePcrawl(w http.ResponseWriter, r *http.Request) {

	u := r.FormValue("url")
	if u == "" {
		http.Error(w, "url required.", http.StatusBadRequest)
		return
	}

	depth, err := strconv.Atoi(r.FormValue("depth"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mxpp, err := strconv.Atoi(r.FormValue("mxpp"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlset := newSyncUrlset()
	ch := pcrawl(u, depth-1, mxpp, urlset)

	b, err := json.Marshal(<-ch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func pcrawl(aUrl string, depth int, mxpp int, urlset *syncUrlset) <-chan *Tree {

	ch := make(chan *Tree)
	go func() {
		ch<-_pcrawl(aUrl, depth, mxpp, urlset)
	}()
	return ch
}

func _pcrawl(aUrl string, depth int, mxpp int, urlset *syncUrlset) *Tree  {

	tree := Tree{Url: aUrl}

	start := time.Now().UnixNano()
	defer func(tree *Tree) {
		tree.Tat = time.Now().UnixNano() - start
	}(&tree)

	resp, err := http.Get(tree.Url)
	if err != nil {
		tree.Error = err.Error()
		return &tree
	}
	defer resp.Body.Close()

	parsedUrl, _ := url.Parse(tree.Url)

	futureChildren := make([]<-chan *Tree, 0, mxpp)

	childNum := 0
	z := html.NewTokenizer(resp.Body)
	for {
		tokenType := z.Next()
		switch tokenType {
		case html.ErrorToken:
			for _, ch := range futureChildren {
				tree.Children = append(tree.Children, <-ch)
			}
			return &tree
		case html.StartTagToken:
			tagname, _ := z.TagName()
			if len(tagname) == 1 && rune(tagname[0]) == 'a' {
				moreAttr := true
				for moreAttr {
					var key, val []byte
					key, val, moreAttr = z.TagAttr()
					log.Printf("href: %s", val)
					var childUrl string
					if string(key) == "href" {

						if string(val[0:4]) == "http" {
							childUrl = string(val)
						} else {
							childUrl = parsedUrl.Scheme + "://" + parsedUrl.Host + string(val)
						}
						log.Printf("childUrl: %s", childUrl)

						if firstHit := func() bool {
							urlset.Lock()
							defer urlset.Unlock()
							if urlset.set[childUrl] {
								return false
							} else {
								urlset.set[childUrl] = true
								return true
							}
						}(); !firstHit {
							continue
						}

						if depth > 0 {
							childCh := pcrawl(childUrl, depth-1, mxpp, urlset)
							futureChildren = append(futureChildren, childCh)
						} else {
							childCh := make(chan *Tree)
							go func() {
								childCh<-&Tree{Url: childUrl}
							}()
							futureChildren = append(futureChildren, childCh)
						}

						childNum++
						if childNum >= mxpp {
							for _, ch := range futureChildren {
								tree.Children = append(tree.Children, <-ch)
							}
							return &tree
						}
						break
					}
				}
			}
		default:
		}
	}
}
