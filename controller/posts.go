package controller

import (
	"germ/model"

	"github.com/gin-gonic/gin"
)

// sharing
type iPostController interface {
	iController
}

// reguires
type iPostRoot interface {
}
type posts struct {
	iPostRoot
}

func (x *posts) register(r *gin.RouterGroup) {
	r.POST("/", x.create)
	r.GET("/:id", x.getOne)
	r.GET("/", x.getAll)
	r.PUT("/", x.update)
	r.DELETE("/:id", x.delete)
}

// Create
func (x *posts) create(c *gin.Context) {

	post, err := model.Post{}.Create(c)
	if err != nil {
		response{false, []string{err.Error()}}.send(c, 400)
		return
	}
	if errs := post.Validate(); errs != nil {
		response{false, errs}.send(c, 400)
		return
	}
	if err := post.Save(); err != nil {
		response{false, []string{err.Error()}}.send(c, 500)
		return
	}
	response{true, nil}.send(c, 201)
}

func (x *posts) getOne(c *gin.Context) {
	post, err := model.Post{}.Create(c.Param("id"))
	if err != nil {
		response{false, []string{err.Error()}}.send(c, 400)
		return
	}
	response{true, post}.send(c, 200)
}
func (x *posts) getAll(c *gin.Context) {
	p, err := model.Post{}.Read()
	if err != nil {
		response{false, []string{err.Error()}}.send(c, 400)
		return
	}
	response{true, p}.send(c, 200)
}

func (x *posts) update(c *gin.Context) {
	post, err := model.Post{}.Create(c)
	if err != nil {
		response{false, []string{err.Error()}}.send(c, 400)
		return
	}
	if err := post.Update(); err != nil {
		response{false, []string{err.Error()}}.send(c, 500)
		return
	}
	response{true, nil}.send(c, 200)
}

func (x *posts) delete(c *gin.Context) {
	post, err := model.Post{}.Create(c.Param(("id")))
	if err != nil {
		response{false, []string{err.Error()}}.send(c, 400)
		return
	}
	if err := post.Destroy(); err != nil {
		response{false, []string{err.Error()}}.send(c, 500)
		return
	}
	response{true, nil}.send(c, 200)
}
