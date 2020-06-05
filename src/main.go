package main

import (
	"github.com/boberneprotiv/notes16/src/crm"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
)

var (
	currentDir, _ = os.Getwd()
	siteFolder    = path.Join(currentDir, "examples", "blog")
	templates     = template.Must(template.New("").
			Funcs(template.FuncMap{
			"escape": func(html string) template.HTML {
				return template.HTML(url.QueryEscape(html))
			},
		}).ParseFiles("templates/components/section-item.html", "templates/components/head.html", "templates/components/navigation.html",
		"templates/pages/dashboard.html", "templates/pages/publications.html", "templates/pages/publication.html", "templates/pages/sections.html"))
)

var sm *crm.SiteManager

func main() {
	manager, err := crm.NewSiteManager(siteFolder)
	sm = manager

	if err != nil {
		log.Fatal(err.Error())
	}

	router := mux.NewRouter().StrictSlash(true)
	// dashboard
	router.HandleFunc("/", dashboardHandler).
		Methods(http.MethodGet)

	contentRouter := router.PathPrefix("/content").Subrouter()

	publicationRouter := contentRouter.PathPrefix("/publication").Subrouter()
	// list of publication
	publicationRouter.HandleFunc("/", publicationListHandler).
		Methods(http.MethodGet)
	// single publication
	publicationRouter.HandleFunc("/{id:.+}", singlePublicationHandler).
		Methods(http.MethodGet)
	// update publication
	publicationRouter.HandleFunc("/{id:.+}/update", updatePublicationHandler).
		Methods(http.MethodPost)

	categoryRouter := contentRouter.PathPrefix("/category").Subrouter()
	// list of categories
	categoryRouter.HandleFunc("/", categoryListHandler).
		Methods(http.MethodGet)
	// single category
	categoryRouter.HandleFunc("/{id:.+}", singleCategoryHandler).
		Methods(http.MethodGet)
	// create category
	categoryRouter.HandleFunc("/create", createCategoryHandler).
		Methods(http.MethodPost)

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "dashboard", nil); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func publicationListHandler(w http.ResponseWriter, r *http.Request) {
	pages := sm.GetSite()
	if err := templates.ExecuteTemplate(w, "index", pages); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func singlePublicationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := url.QueryUnescape(vars["id"])
	page := sm.GetPageById(id)

	if err := templates.ExecuteTemplate(w, "post", page); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func updatePublicationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := url.QueryUnescape(vars["id"])
	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

	c := r.Form["content"][0]
	fm := crm.FrontMatter{
		Title:       r.Form["title"][0],
		Description: r.Form["description"][0],
	}

	page, err := sm.UpdatePageById(id, c, &fm)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

	if err := templates.ExecuteTemplate(w, "post", page); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func categoryListHandler(w http.ResponseWriter, r *http.Request) {
	pages := sm.GetSite()
	if err := templates.ExecuteTemplate(w, "sections", pages); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

	name := r.Form["name"][0]
	p := r.Form["path"][0]

	if err = sm.CreateSection(p, name); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

	pages := sm.GetSite()

	if err := templates.ExecuteTemplate(w, "sections", pages); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func singleCategoryHandler(w http.ResponseWriter, r *http.Request) {
	singlePublicationHandler(w, r)
}
