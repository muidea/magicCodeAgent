package content

// Content core define
type Content struct {
	projectName  string
	rootPath     string
	templatePath string
}

// New new core
func New(rootPath string) *Content {
	return &Content{rootPath: rootPath}
}

// Create create
func (s *Content) Create() (err error) {
	return
}
