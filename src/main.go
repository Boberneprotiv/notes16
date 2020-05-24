package main

import (
	"github.com/boberneprotiv/notes16/src/crm"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
)

var (
	currentDir, _ = os.Getwd()
	siteFolder    = path.Join(currentDir, "examples", "blog")
	templates     = template.Must(template.ParseFiles("templates/publications.html", "templates/publication.html", "templates/head.html"))
)

var sm *crm.SiteManager

func main() {
	manager, err := crm.NewSiteManager(siteFolder)
	sm = manager

	if err != nil {
		log.Fatal(err.Error())
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/post", postHandler)
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		pages := sm.GetSite()
		if err := templates.ExecuteTemplate(w, "index", pages); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getPostHandler(w, r)
	} else {
		postPostHandler(w, r)
	}
}

func getPostHandler(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query()["path"][0]
	page := sm.GetPageByPath(p)

	if err := templates.ExecuteTemplate(w, "post", page); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func postPostHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

	p := r.Form["path"][0]
	c := r.Form["content"][0]
	fm := crm.FrontMatter{
		Title:       r.Form["title"][0],
		Description: r.Form["description"][0],
	}

	page, err := sm.UpdatePageByPath(p, c, &fm)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

	if err := templates.ExecuteTemplate(w, "post", page); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
