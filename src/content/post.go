package content

import (
	"io/ioutil"
	"path"
	"regexp"
)

type Post struct {
	Name        string
	Path        string
	Content     string
	Title       string
	Description string
}

func (s *Site) GetPost(file string) (*Post, error) {
	post, err := openPost(path.Join(s.ContentRoot, file))
	if err != nil {
		return nil, err
	}

	return &Post{
		Name:        "file",
		Path:        file,
		Content:     string(post.Content),
		Title:       post.FrontMatter["title"].(string),
		Description: post.FrontMatter["description"].(string),
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
