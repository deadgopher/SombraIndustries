package controller

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type viewController struct {
	viewRoot
	*template.Template
}

func newViewController(r viewRoot) *viewController {
	x := &viewController{
		r,
		parseTemplates("./view"),
	}
	return x
}

func parseTemplates(p string) *template.Template {
	templ := template.New("")
	err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {

		if strings.Contains(path, ".html") {
			_, err = templ.ParseFiles(path)
			if err != nil {
				log.Println(err)
			}
		}

		return err
	})

	if err != nil {
		panic("---------   AHH IM PANICING!!! ON LINE 48 IN view.go")
	}

	return templ
}

func (x *viewController) showIndex(c *gin.Context) {
	x.ExecuteTemplate(c.Writer, "index", x.Link())
}

func (x *viewController) registerRoutes(r *gin.RouterGroup) {
	r.GET("/", x.showIndex)
}
