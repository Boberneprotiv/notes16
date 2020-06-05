package crm

import (
	"bytes"
	"errors"
	"github.com/gohugoio/hugo/commands"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/parser"
	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/gohugoio/hugo/resources/page"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

type SiteManager struct {
	absPath string
	hugo    *hugolib.HugoSites
}

func NewSiteManager(absPath string) (*SiteManager, error) {
	resp := commands.Execute([]string{"-s", absPath})
	if resp.Err != nil {
		return nil, resp.Err
	}

	sm := &SiteManager{
		absPath: absPath,
	}
	if err := sm.initialize(); err != nil {
		return nil, err
	} else {
		return sm, nil
	}
}

func (s *SiteManager) GetSite() page.Site {
	return s.hugo.Site
}

func (s *SiteManager) UpdatePageById(id string, content string, fm *FrontMatter) (*page.Page, error) {
	p := *s.GetPageById(id)
	if p == nil {
		return nil, errors.New("page not found")
	}

	file, err := os.Open(p.File().Filename())
	if err != nil {
		return nil, err
	}

	pf, err := pageparser.ParseFrontMatterAndContent(file)
	if err != nil {
		return nil, err
	}

	pf.FrontMatter["title"] = fm.Title
	pf.FrontMatter["description"] = fm.Description
	pf.Content = []byte(content)

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
		return nil, err
	}

	newContent.Write(pf.Content)

	if err = ioutil.WriteFile(p.File().Filename(), newContent.Bytes(), 0644); err != nil {
		return nil, err
	}

	if err = s.initialize(); err != nil {
		return nil, err
	}

	return s.GetPageById(id), nil
}

func (s *SiteManager) GetPageById(id string) *page.Page {
	for _, p := range s.hugo.Site.Pages() {
		if p.Path() == id {
			return &p
		}
	}

	return nil
}

func (s *SiteManager) CreateSection(parent string, name string) error {
	p := strings.TrimSuffix(parent, "_index.md")
	resp := commands.Execute([]string{"new", "-s", s.absPath, path.Join(p, name, "_index.md")})
	if resp.Err != nil {
		return resp.Err
	}
	return s.initialize()
}

func (s *SiteManager) initialize() error {
	resp := commands.Execute([]string{"-DEF", "-s", s.absPath})
	if resp.Err != nil {
		return resp.Err
	}

	s.hugo = resp.Result
	return nil
}
