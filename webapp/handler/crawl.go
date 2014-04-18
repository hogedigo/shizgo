package handler

import (
	"encoding/json"
	"log"
	"time"
	"net/http"
	"net/url"
	"strconv"

	"code.google.com/p/go.net/html"
)

type Tree struct {
	Url      string `json:",omitempty"`
	Error    string `json:",omitempty"`
	Children []*Tree
	Tat      int64 `json:",omitempty`
}

func init() {
	http.HandleFunc("/crawl", handleCrawl)
}

func handleCrawl(w http.ResponseWriter, r *http.Request) {

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

	urlset := make(map[string]bool)
	root := crawl(u, depth-1, mxpp, urlset)

	b, err := json.Marshal(root)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func crawl(aUrl string, depth int, mxpp int, urlset map[string]bool) *Tree {

	tree := Tree{Url: aUrl}

	start := time.Now().UnixNano()
	defer func(tree *Tree) {
		tree.Tat = time.Now().UnixNano() - start
	}(&tree);

	resp, err := http.Get(tree.Url)
	if err != nil {
		tree.Error = err.Error()
		return &tree
	}
	defer resp.Body.Close()

	parsedUrl, _ := url.Parse(tree.Url)

	childNum := 0
	z := html.NewTokenizer(resp.Body)
	for {
		tokenType := z.Next()
		switch tokenType {
		case html.ErrorToken:
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

						if urlset[childUrl] {
							continue
						} else {
							urlset[childUrl] = true
						}

						var child *Tree
						if depth > 0 {
							child = crawl(childUrl, depth-1, mxpp, urlset)
						} else {
							child = &Tree{Url: childUrl}
						}
						tree.Children = append(tree.Children, child)

						childNum++
						if childNum >= mxpp {
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
