package controller

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type iTemplateController interface {
	iController
}
type iTemplateRoot interface {
	iRoot
}

type template_controller struct {
	tmpl *template.Template
}

func newTemplateController() *template_controller {
	x := &template_controller{}
	x.parseTemplates("./view")

	return x
}
func (x *template_controller) parseTemplates(y string) {
	templ := template.New("")
	err := filepath.Walk(y, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			x.parseTemplates(path)
		}

		if strings.Contains(path, ".html") {
			_, err = templ.ParseFiles(path)
			if err != nil {
				log.Println(err)
			}
		}

		return err
	})

	if err != nil {
		panic(err)
	}

	x.tmpl = templ
}
