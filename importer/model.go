package importer

type Category string

type Employee struct {
	Name       string
	Categories []Category `json:"interests"`
}

type Gift struct {
	Name       string
	Categories []Category
}
