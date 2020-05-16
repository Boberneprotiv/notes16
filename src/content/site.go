package content

import "path"

type Site struct {
	SiteRoot    string
	ContentRoot string
}

func NewSite(siteFolder string) (*Site, error) {
	return &Site{
		SiteRoot:    siteFolder,
		ContentRoot: path.Join(siteFolder, "content"),
	}, nil
}
