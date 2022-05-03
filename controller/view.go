package controller

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type iViews interface {
	RegisterViews(r *gin.RouterGroup)
	show(gin.ResponseWriter, string)
	setData(interface{})
}

func (x *views) show(w gin.ResponseWriter, to string) {
	x.ExecuteTemplate(w, to, x.iViewData)
}

type views struct {
	viewRoot
	iViewData
	*template.Template
}

func newViewController(r viewRoot) *views {
	x := &views{}
	x.viewRoot = r
	x.iViewData = viewData{esiLink: x.esiLink()}.new()
	x.Template = parseTemplates("./view")
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

func (x *views) index(c *gin.Context) {
	x.ExecuteTemplate(c.Writer, "index", x.iViewData)
}

func (x *views) RegisterViews(r *gin.RouterGroup) {
	r.GET("/", x.index)
}
