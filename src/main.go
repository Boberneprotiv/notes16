package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"regexp"
	"strings"
)

var (
	siteFolder    = "XXX"
	contentFolder = path.Join(siteFolder, "content")
	templates     = template.Must(template.ParseFiles("templates/index.html", "templates/post.html", "templates/category-list.html"))
)

type File struct {
	Name      string
	Path      string
	IsCatalog bool
	Files     []File
}

type Post struct {
	Name    string
	Path    string
	Content string
}

type Catalog struct {
	Name  string
	Files []File
}

func NewCatalog(name string, path string, files []File) File {
	return File{
		Name:      name,
		Path:      path,
		IsCatalog: true,
		Files:     files,
	}
}

func NewFile(name string, path string) File {
	return File{
		Name:      name,
		Path:      path,
		IsCatalog: false,
		Files:     make([]File, 0),
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/post", postHandler)
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func contentCatalog() Catalog {
	return Catalog{
		Name: "Hugo catalog",
		Files: []File{
			NewCatalog("Home page", "_index.md", discoverFolder(contentFolder)),
		},
	}
}

func discoverFolder(folder string) []File {
	files, _ := ioutil.ReadDir(folder)
	result := make([]File, 0)
	for _, file := range files {
		p := path.Join(folder, file.Name())
		if file.IsDir() {
			result = append(result, NewCatalog(file.Name(), path.Join(strings.TrimPrefix(p, contentFolder), "_index.md"), discoverFolder(p)))
		} else if r, _ := regexp.MatchString(`(?m)_index(\.ua)?\.md`, file.Name()); !r {
			result = append(result, NewFile(file.Name(), strings.TrimPrefix(p, contentFolder)))
		}
	}

	return result
}

func getPost(file string) Post {
	path := path.Join(contentFolder, file)
	b, _ := ioutil.ReadFile(path)
	content := string(b)

	re := regexp.MustCompile(`(?s)---(.*?)---`)

	return Post{
		Name:    "file",
		Path:    file,
		Content: strings.TrimPrefix(content, re.FindStringSubmatch(content)[0]),
	}
}

func updateContent(file string, newContent string) {
	path := path.Join(contentFolder, file)
	b, _ := ioutil.ReadFile(path)
	re := regexp.MustCompile(`(?s)---(.*?)---`)
	rawFile := string(b)

	metaInfo := re.FindStringSubmatch(rawFile)[0]

	ioutil.WriteFile(path, []byte(metaInfo+"\n"+newContent), 0644)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "index", contentCatalog()); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		file := r.URL.Query()["file"][0]
		if err := templates.ExecuteTemplate(w, "post", getPost(file)); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	} else {
		if err := r.ParseForm(); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}

		file := r.Form.Get("path")
		content := r.Form.Get("content")

		updateContent(file, content)

		if err := templates.ExecuteTemplate(w, "post", getPost(file)); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	}
}
