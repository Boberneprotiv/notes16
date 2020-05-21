package main

import (
	"github.com/gohugoio/hugo/commands"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/resources/page"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
)

var (
	currentDir, _ = os.Getwd()
	siteFolder    = path.Join(currentDir, "examples", "blog")
	templates     = template.Must(template.ParseFiles("templates/index.html", "templates/post.html", "templates/head.html", "templates/category-list.html"))
)

var s *hugolib.HugoSites

func main() {
	resp := commands.Execute([]string{"-s", siteFolder})
	if resp.Err != nil {
		log.Fatal(resp.Err.Error())
	}
	s = resp.Result

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/post", postHandler)
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var home page.Page
		pages := s.Site.Pages()
		for _, p := range pages {
			if p.IsHome() {
				home = p
			}
		}
		if err := templates.ExecuteTemplate(w, "index", home); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	getPostHandler(w, r)
}

func getPostHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query()["path"][0]
	var page page.Page
	pages := s.Site.Pages()
	for _, p := range pages {
		if p.Path() == path {
			page = p
		}
	}

	if err := templates.ExecuteTemplate(w, "post", page); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
