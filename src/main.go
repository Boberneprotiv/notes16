package main

import (
	"bytes"
	"github.com/gohugoio/hugo/commands"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/parser"
	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/gohugoio/hugo/resources/page"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

var (
	currentDir, _ = os.Getwd()
	siteFolder    = path.Join(currentDir, "examples", "blog")
	templates     = template.Must(template.ParseFiles("templates/index.html", "templates/post.html", "templates/head.html", "templates/category-list.html"))
)

var s *hugolib.HugoSites

func reinit() {
	resp := commands.Execute([]string{"-s", siteFolder})
	if resp.Err != nil {
		log.Fatal(resp.Err.Error())
	}
	s = resp.Result
}

func main() {
	reinit()

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
	if r.Method == http.MethodGet {
		getPostHandler(w, r)
	} else {
		postPostHandler(w, r)
	}
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

func postPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	path := r.Form["path"][0]
	var page page.Page
	pages := s.Site.Pages()
	for _, p := range pages {
		if p.Path() == path {
			page = p
		}
	}

	c := &UpdatePageCommand{
		Title:       r.Form["title"][0],
		Description: r.Form["description"][0],
		Content:     r.Form["content"][0],
	}

	err := overwritePage(c, page)
	reinit()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

	if err := templates.ExecuteTemplate(w, "post", page); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

type UpdatePageCommand struct {
	Title       string
	Description string
	Content     string
}

func overwritePage(c *UpdatePageCommand, p page.Page) error {
	file, err := os.Open(p.File().Filename())
	if err != nil {
		return err
	}

	pf, err := pageparser.ParseFrontMatterAndContent(file)
	if err != nil {
		return err
	}

	pf.FrontMatter["title"] = c.Title
	pf.FrontMatter["description"] = c.Description
	pf.Content = []byte(c.Content)

	if pf.FrontMatterFormat == metadecoders.JSON || pf.FrontMatterFormat == metadecoders.YAML || pf.FrontMatterFormat == metadecoders.TOML {
		for k, v := range pf.FrontMatter {
			switch vv := v.(type) {
			case time.Time:
				pf.FrontMatter[k] = vv.Format(time.RFC3339)
			}
		}
	}

	var newContent bytes.Buffer
	err = parser.InterfaceToFrontMatter(pf.FrontMatter, metadecoders.YAML, &newContent)
	if err != nil {
		return err
	}

	newContent.Write(pf.Content)

	ioutil.WriteFile(p.File().Filename(), newContent.Bytes(), 0644)

	return nil
}
