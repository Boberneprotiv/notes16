package content

type File struct {
	Name      string
	Path      string
	IsCatalog bool
	Files     []File
}
