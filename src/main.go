package main

import (
	"github.com/boberneprotiv/notes16/src/content"
	"html/template"
	"log"
	"net/http"
)

var (
	siteFolder = "/Users/vladimirborodin/Documents/Avtodoctor/HugoSite"
	site, _    = content.NewSite(siteFolder)
	templates  = template.Must(template.ParseFiles("templates/index.html", "templates/post.html", "templates/category-list.html"))
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/post", postHandler)
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if c, err := site.GetCatalog(); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		} else if err := templates.ExecuteTemplate(w, "index", c); err != nil {
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
	file := r.URL.Query()["file"][0]

	if p, err := site.GetPost(file); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	} else if err := templates.ExecuteTemplate(w, "post", p); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func postPostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

	file := r.Form.Get("path")
	c := r.Form.Get("content")

	if err := site.UpdatePost(file, c); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

	if p, err := site.GetPost(file); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	} else if err := templates.ExecuteTemplate(w, "post", p); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
