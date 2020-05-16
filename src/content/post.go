package content

import (
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

type Post struct {
	Name    string
	Path    string
	Content string
}

func (s *Site) GetPost(file string) (*Post, error) {
	p := path.Join(s.ContentRoot, file)
	b, _ := ioutil.ReadFile(p)
	content := string(b)

	re := regexp.MustCompile(`(?s)---(.*?)---`)

	return &Post{
		Name:    "file",
		Path:    file,
		Content: strings.TrimPrefix(content, re.FindStringSubmatch(content)[0]),
	}, nil
}

func (s *Site) UpdatePost(file string, newContent string) error {
	p := path.Join(s.ContentRoot, file)
	b, _ := ioutil.ReadFile(p)

	re := regexp.MustCompile(`(?s)---(.*?)---`)
	rawFile := string(b)

	metaInfo := re.FindStringSubmatch(rawFile)[0]

	return ioutil.WriteFile(p, []byte(metaInfo+"\n"+newContent), 0644)
}
