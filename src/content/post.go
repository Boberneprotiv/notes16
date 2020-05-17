package content

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

type Post struct {
	Name        string
	Path        string
	Content     string
	Title       string
	Description string
}

func (s *Site) GetPost(file string) (*Post, error) {
	p := path.Join(s.ContentRoot, file)
	b, _ := ioutil.ReadFile(p)
	content := string(b)

	re := regexp.MustCompile(`(?s)---(.*?)---`)
	rawMeta := re.FindStringSubmatch(content)[0]

	var meta map[string]interface{}
	if err := yaml.Unmarshal([]byte(strings.Trim(rawMeta, "---")), &meta); err != nil {
		return nil, err
	}

	return &Post{
		Name:        "file",
		Path:        file,
		Content:     strings.TrimPrefix(content, rawMeta),
		Title:       meta["title"].(string),
		Description: meta["description"].(string),
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
