package content

type Catalog struct {
	Name  string
	Files []File
}

func (s *Site) GetCatalog() (*Catalog, error) {
	return &Catalog{
		Name: "Hugo catalog",
		Files: []File{
			newCatalog("Home page", "_index.md", discoverFolder(s.ContentRoot, s.ContentRoot)),
		},
	}, nil
}
