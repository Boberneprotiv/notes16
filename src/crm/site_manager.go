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

func (s *SiteManager) UpdatePageByPath(path string, content string, fm *FrontMatter) (*page.Page, error) {
	p := *s.GetPageByPath(path)
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

	return s.GetPageByPath(path), nil
}

func (s *SiteManager) GetPageByPath(path string) *page.Page {
	for _, p := range s.hugo.Site.Pages() {
		if p.Path() == path {
			return &p
		}
	}

	return nil
}

func (s *SiteManager) initialize() error {
	resp := commands.Execute([]string{"-s", s.absPath})
	if resp.Err != nil {
		return resp.Err
	}

	s.hugo = resp.Result
	return nil
}
