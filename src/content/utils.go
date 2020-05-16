package content

import (
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

func newCatalog(name string, path string, files []File) File {
	return File{
		Name:      name,
		Path:      path,
		IsCatalog: true,
		Files:     files,
	}
}

func newFile(name string, path string) File {
	return File{
		Name:      name,
		Path:      path,
		IsCatalog: false,
		Files:     make([]File, 0),
	}
}

func discoverFolder(contentRoot string, folder string) []File {
	files, _ := ioutil.ReadDir(folder)
	result := make([]File, 0)
	for _, file := range files {
		p := path.Join(folder, file.Name())
		relName := relativeFileName(contentRoot, p)
		if file.IsDir() {
			result = append(result, newCatalog(file.Name(), path.Join(relName, "_index.md"), discoverFolder(contentRoot, p)))
		} else if !isCategory(file.Name()) {
			result = append(result, newFile(file.Name(), relName))
		}
	}

	return result
}

func relativeFileName(contentRoot string, absFN string) string {
	return strings.TrimPrefix(absFN, contentRoot)
}

func isCategory(fileName string) bool {
	r, _ := regexp.MatchString(`(?m)_index(\.ua)?\.md`, fileName)
	return r
}
